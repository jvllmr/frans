package apiRoutes

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent"
	"github.com/jvllmr/frans/pkg/ent/ticket"
	"github.com/jvllmr/frans/pkg/util"
)

type PublicTicket struct {
	Id              uuid.UUID    `json:"id"`
	Comment         *string      `json:"comment"`
	EstimatedExpiry string       `json:"estimatedExpiry"`
	User            PublicUser   `json:"owner"`
	Files           []PublicFile `json:"files"`
	CreatedAt       string       `json:"createdAt"`
}

func ToPublicTicket(configValue config.Config, ticket *ent.Ticket) PublicTicket {
	files := []PublicFile{}
	for _, file := range ticket.Edges.Files {
		files = append(files, ToPublicFile(file))
	}

	return PublicTicket{
		Id:              ticket.ID,
		Comment:         ticket.Comment,
		User:            ToPublicUser(ticket.Edges.Owner),
		EstimatedExpiry: util.GetEstimatedExpiry(configValue, ticket).Format(http.TimeFormat),
		Files:           files,
		CreatedAt:       ticket.CreatedAt.UTC().Format(http.TimeFormat),
	}
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
}

func createTicketFactory(configValue config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet(config.UserGinContext).(*ent.User)
		var form ticketForm
		tx, err := config.DBClient.Tx(c.Request.Context())
		if err != nil {
			c.AbortWithError(500, err)
		}
		if err := c.ShouldBind(&form); err == nil {
			salt, err := util.GenerateSalt()
			if err != nil {
				c.AbortWithError(500, err)
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
				SetOwner(user)

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
			if len(files) > int(configValue.MaxFiles) {
				c.AbortWithStatus(http.StatusUnprocessableEntity)
				return
			}
			util.EnsureFilesTmpPath(configValue)

			for _, fileHeader := range files {
				if fileHeader.Size > configValue.MaxSizes {
					c.AbortWithStatus(http.StatusUnprocessableEntity)
					return
				}

				incomingFileHandle, _ := fileHeader.Open()
				hasher := sha512.New()
				tmpFilePath := util.GetFilesTmpFilePath(configValue)
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
					Save(c.Request.Context())
				if err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
				}
				targetFilePath := util.GetFilesFilePath(configValue, sha512sum)
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

			c.JSON(http.StatusCreated, ToPublicTicket(configValue, ticketValue))
			tx.Commit()
		} else {
			c.AbortWithError(422, err)
		}

	}
}

func setupTicketGroup(r *gin.RouterGroup, configValue config.Config) {
	r.POST("", createTicketFactory(configValue))
}
