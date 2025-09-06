package apiRoutes

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent/ticket"
	"github.com/jvllmr/frans/internal/ent/user"
	"github.com/jvllmr/frans/internal/mail"
	"github.com/jvllmr/frans/internal/middleware"
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/util"
)

type ticketController struct {
	config        config.Config
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
	currentUser := middleware.GetCurrentUser(c)
	var form ticketForm
	tx, err := config.DBClient.Tx(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	if err := c.ShouldBind(&form); err == nil {
		salt := util.GenerateSalt()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
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

		ticketValue, err := ticketBuilder.Save(c.Request.Context())
		if err != nil {
			c.AbortWithError(400, err)
		}

		multipartForm, _ := c.MultipartForm()
		files := multipartForm.File["files[]"]
		if len(files) > int(tc.config.MaxFiles) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		tc.fileService.EnsureFilesTmpPath()

		for _, fileHeader := range files {
			if fileHeader.Size > tc.config.MaxSizes {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			incomingFileHandle, _ := fileHeader.Open()
			hasher := sha512.New()
			tmpFilePath := tc.fileService.FilesTmpFilePath()
			tmpFileHandle, err := os.Create(tmpFilePath)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			defer os.Remove(tmpFilePath)
			writer := io.MultiWriter(hasher, tmpFileHandle)
			_, err = io.Copy(writer, incomingFileHandle)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			tmpFileHandle.Close()

			hash := hasher.Sum(nil)
			sha512sum := hex.EncodeToString(hash)
			dbFile, err := tx.File.Create().
				SetID(uuid.New()).
				SetName(fileHeader.Filename).
				SetSize(uint64(fileHeader.Size)).
				SetSha512(sha512sum).
				SetExpiryType(form.ExpiryType).
				SetExpiryDaysSinceLastDownload(form.ExpiryDaysSinceLastDownload).
				SetExpiryTotalDays(form.ExpiryTotalDays).
				SetExpiryTotalDownloads(form.ExpiryTotalDownloads).
				Save(c.Request.Context())
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			targetFilePath := tc.fileService.FilesFilePath(sha512sum)
			if _, err = os.Stat(targetFilePath); err != nil {
				os.Rename(tmpFilePath, targetFilePath)
			}
			ticketValue = tx.Ticket.UpdateOne(ticketValue).
				AddFiles(dbFile).
				SaveX(c.Request.Context())
		}

		ticketValue = tx.Ticket.Query().
			Where(ticket.ID(ticketValue.ID)).
			WithFiles().
			WithOwner().
			OnlyX(c.Request.Context())

		err = util.RefreshUserTotalDataSize(c.Request.Context(), currentUser)
		if err != nil {
			slog.Error(
				"Could not refresh total data size of user",
				"err",
				err,
				"user",
				currentUser.Username,
			)
		}
		tx.User.UpdateOne(currentUser).AddSubmittedTickets(1).SaveX(c.Request.Context())
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
		c.AbortWithError(422, err)
	}

}

func (tc *ticketController) fetchTicketsHandler(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	query := config.DBClient.Ticket.Query().WithFiles().WithOwner()

	if !currentUser.IsAdmin {
		query = query.Where(ticket.HasOwnerWith(user.ID(currentUser.ID)))
	}

	tickets, err := query.All(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	publicTickets := make([]services.PublicTicket, len(tickets))
	for i, ticketValue := range tickets {
		publicTickets[i] = tc.ticketService.ToPublicTicket(tc.fileService, ticketValue)
	}
	c.JSON(http.StatusOK, publicTickets)
}

func setupTicketGroup(r *gin.RouterGroup, configValue config.Config) {
	controller := ticketController{
		config:        configValue,
		ticketService: services.NewTicketService(configValue),
		fileService:   services.NewFileService(configValue),
		mailer:        mail.NewMailer(configValue),
	}
	r.POST("", controller.createTicketHandler)
	r.GET("", controller.fetchTicketsHandler)
}
