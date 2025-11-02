package shareRoutes

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/file"
	"github.com/jvllmr/frans/internal/ent/ticket"
	"github.com/jvllmr/frans/internal/mail"
	"github.com/jvllmr/frans/internal/otel"
	apiTypes "github.com/jvllmr/frans/internal/routes/api/types"
	"github.com/jvllmr/frans/internal/services"

	"github.com/jvllmr/frans/internal/util"
)

type ticketShareController struct {
	config        config.Config
	db            *ent.Client
	ticketService services.TicketService
	fileService   services.FileService
	mailer        mail.Mailer
}

func (tsc *ticketShareController) fetchTicket(c *gin.Context) {
	_, span := otel.NewSpan(c.Request.Context(), "fetchTicketShare")
	defer span.End()
	ticketValue := c.MustGet(config.ShareTicketContext).(*ent.Ticket)
	c.JSON(http.StatusOK, tsc.ticketService.ToPublicTicket(tsc.fileService, ticketValue))
}

func (tsc *ticketShareController) fetchTicketAccessToken(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchTicketShareAccessToken")
	defer span.End()
	tokenValueBytes := util.GenerateSalt()

	tokenValue := hex.EncodeToString(tokenValueBytes)
	token := tsc.db.ShareAccessToken.Create().
		SetID(tokenValue).
		SetExpiry(time.Now().Add(10 * time.Second)).
		SaveX(ctx)
	c.SetCookie(
		config.ShareAccessTokenCookieName,
		token.ID,
		10,
		strings.TrimSuffix(c.Request.URL.Path, "/token"),
		"",
		false,
		true,
	)
	c.JSON(http.StatusCreated, apiTypes.PublicShareAccessToken{Token: token.ID})
}

func (tsc *ticketShareController) fetchTicketFile(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchTicketShareFile")
	defer span.End()
	var requestedFile apiTypes.RequestedFileParam
	if err := c.ShouldBindUri(&requestedFile); err != nil {
		util.GinAbortWithError(c, http.StatusBadRequest, err)
		return
	}
	fileValue, err := tsc.db.File.Query().
		WithData().
		Where(file.ID(uuid.MustParse(requestedFile.ID))).
		Only(ctx)
	if err != nil {
		util.GinAbortWithError(c, http.StatusNotFound, err)
	}
	_, err = tsc.db.File.UpdateOne(fileValue).
		SetLastDownload(time.Now()).
		AddTimesDownloaded(1).
		Save(ctx)
	if err != nil {
		util.GinAbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	filePath := tsc.fileService.FilesFilePath(fileValue.Edges.Data.ID)
	c.FileAttachment(filePath, fileValue.Name)
	ticketValue := c.MustGet(config.ShareTicketContext).(*ent.Ticket)
	if ticketValue.EmailOnDownload != nil &&
		(fileValue.LastDownload == nil || fileValue.LastDownload.Before(ticketValue.CreatedAt)) {
		tsc.mailer.SendFileDownloadNotification(
			c,
			*ticketValue.EmailOnDownload,
			ticketValue,
			fileValue,
		)
	}
}

func setupTicketShareRoutes(r *gin.RouterGroup, configValue config.Config, db *ent.Client) {
	getTicketMiddleware := func(c *gin.Context) {
		ctx, span := otel.NewSpan(c.Request.Context(), "checkTicketShareAuth")
		defer span.End()
		ticketId := c.Param("ticketId")
		username, password, ok := c.Request.BasicAuth()

		if !ok {
			tokenCookie, err := c.Cookie(config.ShareAccessTokenCookieName)
			if err != nil {
				util.GinAbortWithError(c, http.StatusUnauthorized, err)
				return
			}
			token, err := db.ShareAccessToken.Get(ctx, tokenCookie)
			if err != nil {
				util.GinAbortWithError(c, http.StatusUnauthorized, err)
			} else if token.Expiry.Before(time.Now()) {
				util.GinAbortWithError(c, http.StatusUnauthorized, fmt.Errorf("token expired"))
				return
			}
		} else if username != ticketId {
			util.GinAbortWithError(c, http.StatusUnauthorized, fmt.Errorf("token does not match share"))
			return
		}
		uuidValue, err := uuid.Parse(ticketId)
		if err != nil {
			util.GinAbortWithError(c, http.StatusBadRequest, err)
		}
		ticketValue, err := db.Ticket.Query().
			Where(ticket.ID(uuidValue)).
			WithOwner().
			WithFiles(func(fq *ent.FileQuery) { fq.WithData() }).Only(ctx)

		if err != nil {
			util.GinAbortWithError(c, http.StatusUnauthorized, err)
			return
		}

		if ok && !util.VerifyPassword(password, ticketValue.HashedPassword, ticketValue.Salt) {
			util.GinAbortWithError(c, http.StatusUnauthorized, fmt.Errorf("password incorrect"))
		}

		c.Set(config.ShareTicketContext, ticketValue)
	}

	singleTicketShareGroup := r.Group("/:ticketId", getTicketMiddleware)
	controller := ticketShareController{
		config:        configValue,
		db:            db,
		ticketService: services.NewTicketService(configValue),
		fileService:   services.NewFileService(configValue, db),
		mailer:        mail.NewMailer(configValue),
	}

	singleTicketShareGroup.GET("", controller.fetchTicket)

	singleTicketShareGroup.GET("/token", controller.fetchTicketAccessToken)

	singleTicketShareGroup.GET("/file/:fileId", controller.fetchTicketFile)

}
