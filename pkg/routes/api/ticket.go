package apiRoutes

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent"
	"github.com/jvllmr/frans/pkg/util"
)

type PublicTicket struct {
	Comment         *string      `json:"comment"`
	EstimatedExpiry string       `json:"estimatedExpiry"`
	User            PublicUser   `json:"owner"`
	Files           []PublicFile `json:"files"`
}

func ToPublicTicket(configValue config.Config, ticket *ent.Ticket) PublicTicket {
	files := []PublicFile{}
	for _, file := range files {
		files = append(files, file)
	}

	return PublicTicket{
		Comment:         ticket.Comment,
		User:            ToPublicUser(ticket.Edges.Owner),
		EstimatedExpiry: util.GetEstimatedExpiry(configValue, ticket).Format(http.TimeFormat),
		Files:           files,
	}
}

type ticketForm struct {
	Comment                     *string `binding:"required" form:"comment"`
	Email                       *string `binding:"required" form:"email"`
	Password                    string  `binding:"required" form:"password"`
	EmailPassword               bool    `binding:"required" form:"emailPassword"`
	ExpiryType                  string  `binding:"required" form:"expiryType"`
	ExpiryTotalDays             uint8   `binding:"required" form:"expiryTotalDays"`
	ExpiryDaysSinceLastDownload uint8   `binding:"required" form:"expiryDaysSinceLastDownload"`
	ExpiryTotalDownloads        uint8   `binding:"required" form:"expiryTotalDownloads"`
	EmailOnDownload             *string `binding:"required" form:"emailOnDownload"`
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
			ticket, err := tx.Ticket.Create().
				SetComment(*form.Comment).
				SetEmailOnDownload(*form.EmailOnDownload).
				SetExpiryType(form.ExpiryType).
				SetExpiryDaysSinceLastDownload(form.ExpiryDaysSinceLastDownload).
				SetExpiryTotalDays(form.ExpiryTotalDays).
				SetExpiryTotalDownloads(form.ExpiryTotalDownloads).
				SetHashedPassword(hashedPassword).
				SetSalt(hex.EncodeToString(salt)).
				SetOwner(user).
				Save(c.Request.Context())
			if err != nil {
				c.AbortWithError(400, err)
			}

			multipartForm, _ := c.MultipartForm()
			files := multipartForm.File["files"]
			for _, fileHeader := range files {

				incomingFileHandle, _ := fileHeader.Open()
				hasher := sha256.New()
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
				sha256sum := hex.EncodeToString(hash)
				dbFile, err := tx.File.Get(c.Request.Context(), sha256sum)
				if err != nil {
					if ent.IsNotFound(err) {
						dbFile, err = tx.File.Create().
							SetID(sha256sum).
							SetName(fileHeader.Filename).
							SetSize(uint64(fileHeader.Size)).
							Save(c.Request.Context())
						if err != nil {
							c.AbortWithError(http.StatusInternalServerError, err)
						} else {
							os.Rename(tmpFilePath, util.GetFilesFilePath(configValue, sha256sum))
						}
					} else {
						c.AbortWithError(http.StatusInternalServerError, err)
					}
				}
				ticket = config.DBClient.Ticket.UpdateOne(ticket).
					AddFiles(dbFile).
					SaveX(c.Request.Context())
			}
			tx.Commit()
			c.JSON(http.StatusCreated, ToPublicTicket(configValue, ticket))
		} else {
			c.AbortWithError(422, err)
		}

	}
}

func setupTicketGroup(r *gin.RouterGroup, configValue config.Config) {
	r.POST("", createTicketFactory(configValue))
}
