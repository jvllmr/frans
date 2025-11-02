package apiRoutes

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/ticket"
	"github.com/jvllmr/frans/internal/ent/user"
	"github.com/jvllmr/frans/internal/mail"
	"github.com/jvllmr/frans/internal/middleware"
	"github.com/jvllmr/frans/internal/otel"
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/util"
)

type ticketController struct {
	config        config.Config
	db            *ent.Client
	ticketService services.TicketService
	fileService   services.FileService
	mailer        mail.Mailer
}

type ticketForm struct {
	Comment                     *string `form:"comment"`
	Email                       *string `form:"email"`
	Password                    string  `form:"password"                    binding:"required"`
	EmailPassword               bool    `form:"emailPassword"`
	ExpiryType                  string  `form:"expiryType"                  binding:"required"`
	ExpiryTotalDays             uint8   `form:"expiryTotalDays"             binding:"required"`
	ExpiryDaysSinceLastDownload uint8   `form:"expiryDaysSinceLastDownload" binding:"required"`
	ExpiryTotalDownloads        uint8   `form:"expiryTotalDownloads"        binding:"required"`
	EmailOnDownload             *string `form:"emailOnDownload"`
	CreatorLang                 string  `form:"creatorLang"                 binding:"required"`
	ReceiverLang                string  `form:"receiverLang"                binding:"required"`
}

func (tc *ticketController) createTicketHandler(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "createTicket")
	defer span.End()

	currentUser := middleware.GetCurrentUser(c)
	var form ticketForm
	tx, err := tc.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		util.GinAbortWithError(c, http.StatusInternalServerError, err)
	}
	if err := c.ShouldBind(&form); err == nil {
		salt := util.GenerateSalt()

		hashedPassword := util.HashPassword(form.Password, salt)
		ticketBuilder := tx.Ticket.Create().
			SetID(uuid.New()).
			SetExpiryType(form.ExpiryType).
			SetExpiryDaysSinceLastDownload(form.ExpiryDaysSinceLastDownload).
			SetExpiryTotalDays(form.ExpiryTotalDays).
			SetExpiryTotalDownloads(form.ExpiryTotalDownloads).
			SetHashedPassword(hashedPassword).
			SetSalt(hex.EncodeToString(salt)).
			SetOwner(currentUser).
			SetCreatorLang(form.CreatorLang)

		if form.Comment != nil {
			ticketBuilder = ticketBuilder.SetComment(*form.Comment)
		}

		if form.EmailOnDownload != nil {
			ticketBuilder = ticketBuilder.SetEmailOnDownload(*form.EmailOnDownload)
		}

		ticketValue, err := ticketBuilder.Save(ctx)
		if err != nil {
			util.GinAbortWithError(c, http.StatusBadRequest, err)
		}

		multipartForm, _ := c.MultipartForm()
		files := multipartForm.File["files[]"]
		if len(files) > int(tc.config.MaxFiles) {
			util.GinAbortWithError(
				c,
				http.StatusBadRequest,
				fmt.Errorf(
					"maximum of %d files allowed per upload. %d uploaded",
					tc.config.MaxFiles,
					len(files),
				),
			)
			return
		}
		tc.fileService.EnsureFilesTmpPath()

		for _, fileHeader := range files {
			dbFile, err := tc.fileService.CreateFile(
				ctx,
				tx, fileHeader, currentUser,
				ticketValue.ExpiryType,
				ticketValue.ExpiryDaysSinceLastDownload,
				ticketValue.ExpiryTotalDays,
				ticketValue.ExpiryTotalDownloads,
			)
			if err != nil {
				var errFileTooBig *services.ErrFileTooBig
				if errors.As(err, &errFileTooBig) {
					util.GinAbortWithError(c, http.StatusBadRequest, err)
				} else {
					util.GinAbortWithError(c, http.StatusInternalServerError, err)
				}
				return
			}
			ticketValue = tx.Ticket.UpdateOne(ticketValue).
				AddFiles(dbFile).
				SaveX(ctx)
		}

		ticketValue = tx.Ticket.Query().
			Where(ticket.ID(ticketValue.ID)).
			WithFiles(func(fq *ent.FileQuery) { fq.WithData() }).
			WithOwner().
			OnlyX(ctx)

		err = util.RefreshUserTotalDataSize(ctx, currentUser, tx)
		if err != nil {
			slog.Error(
				"Could not refresh total data size of user",
				"err",
				err,
				"user",
				currentUser.Username,
			)
		}
		tx.User.UpdateOne(currentUser).AddSubmittedTickets(1).SaveX(ctx)
		c.JSON(http.StatusCreated, tc.ticketService.ToPublicTicket(tc.fileService, ticketValue))
		if form.Email != nil {
			var toBeEmailedPassword *string = nil
			if form.EmailPassword {
				toBeEmailedPassword = &form.Password
			}
			tc.mailer.SendTicketSharedNotification(
				c,
				tc.ticketService,
				*form.Email,
				form.ReceiverLang,
				ticketValue,
				toBeEmailedPassword,
			)
		}

		tx.Commit()
	} else {
		util.GinAbortWithError(c, http.StatusUnprocessableEntity, err)
	}

}

func (tc *ticketController) fetchTicketsHandler(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchTickets")
	defer span.End()
	currentUser := middleware.GetCurrentUser(c)
	query := tc.db.Ticket.Query().WithFiles(func(fq *ent.FileQuery) { fq.WithData() }).WithOwner()

	if !currentUser.IsAdmin {
		query = query.Where(ticket.HasOwnerWith(user.ID(currentUser.ID)))
	}

	tickets, err := query.All(ctx)
	if err != nil {
		util.GinAbortWithError(c, http.StatusInternalServerError, err)
	}
	publicTickets := make([]services.PublicTicket, len(tickets))
	for i, ticketValue := range tickets {
		publicTickets[i] = tc.ticketService.ToPublicTicket(tc.fileService, ticketValue)
	}
	c.JSON(http.StatusOK, publicTickets)
}

func setupTicketGroup(r *gin.RouterGroup, configValue config.Config, db *ent.Client) {
	controller := ticketController{
		config:        configValue,
		db:            db,
		ticketService: services.NewTicketService(configValue),
		fileService:   services.NewFileService(configValue, db),
		mailer:        mail.NewMailer(configValue),
	}
	r.POST("", controller.createTicketHandler)
	r.GET("", controller.fetchTicketsHandler)
}
