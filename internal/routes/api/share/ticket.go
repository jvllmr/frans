package shareRoutes

import (
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/ticket"
	"github.com/jvllmr/frans/internal/mail"
	apiTypes "github.com/jvllmr/frans/internal/routes/api/types"

	"github.com/jvllmr/frans/internal/util"
)

func getTicketMiddleware(c *gin.Context) {

	ticketId := c.Param("ticketId")
	username, password, ok := c.Request.BasicAuth()

	if !ok {
		tokenCookie, err := c.Cookie(config.ShareAccessTokenCookieName)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		token, err := config.DBClient.ShareAccessToken.Get(c.Request.Context(), tokenCookie)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		} else if token.Expiry.Before(time.Now()) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	} else if username != ticketId {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	uuidValue, err := uuid.Parse(ticketId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	ticketValue, err := config.DBClient.Ticket.Query().
		Where(ticket.ID(uuidValue)).
		WithOwner().
		WithFiles().Only(c.Request.Context())

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if ok && !util.VerifyPassword(password, ticketValue.HashedPassword, ticketValue.Salt) {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Set(config.ShareTicketContext, ticketValue)
}

func setupTicketShareRoutes(r *gin.RouterGroup, configValue config.Config) {
	singleTicketShareGroup := r.Group("/:ticketId", getTicketMiddleware)

	singleTicketShareGroup.GET("", func(c *gin.Context) {
		ticketValue := c.MustGet(config.ShareTicketContext).(*ent.Ticket)
		c.JSON(http.StatusOK, apiTypes.ToPublicTicket(configValue, ticketValue))
	})

	singleTicketShareGroup.GET("/token", func(c *gin.Context) {
		tokenValueBytes, err := util.GenerateSalt()
		if err != nil {
			panic(err)
		}
		tokenValue := hex.EncodeToString(tokenValueBytes)
		token := config.DBClient.ShareAccessToken.Create().
			SetID(tokenValue).
			SetExpiry(time.Now().Add(10 * time.Second)).
			SaveX(c.Request.Context())
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
	})

	singleTicketShareGroup.GET("/file/:fileId", func(c *gin.Context) {
		var requestedFile apiTypes.RequestedFileParam
		if err := c.ShouldBindUri(&requestedFile); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		fileValue, err := config.DBClient.File.Get(
			c.Request.Context(),
			uuid.MustParse(requestedFile.ID),
		)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}
		fileValue = config.DBClient.File.UpdateOne(fileValue).
			SetLastDownload(time.Now()).
			AddTimesDownloaded(1).
			SaveX(c.Request.Context())
		filePath := util.FilesFilePath(configValue, fileValue.Sha512)
		c.FileAttachment(filePath, fileValue.Name)
		ticketValue := c.MustGet(config.ShareTicketContext).(*ent.Ticket)
		if ticketValue.EmailOnDownload != nil &&
			(fileValue.LastDownload == nil || fileValue.LastDownload.Before(ticketValue.CreatedAt)) {
			mail.SendFileDownloadNotification(
				c,
				configValue,
				*ticketValue.EmailOnDownload,
				ticketValue,
				fileValue,
			)
		}
	})

}
