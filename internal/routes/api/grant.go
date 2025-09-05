package apiRoutes

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/grant"
	"github.com/jvllmr/frans/internal/ent/user"
	"github.com/jvllmr/frans/internal/mail"
	"github.com/jvllmr/frans/internal/middleware"
	apiTypes "github.com/jvllmr/frans/internal/routes/api/types"
	"github.com/jvllmr/frans/internal/util"
)

type grantForm struct {
	Comment                         *string `form:"comment"`
	Email                           *string `form:"email"`
	Password                        string  `form:"password"                        binding:"required"`
	EmailPassword                   bool    `form:"emailPassword"`
	ExpiryType                      string  `form:"expiryType"                      binding:"required"`
	ExpiryTotalDays                 uint8   `form:"expiryTotalDays"                 binding:"required"`
	ExpiryDaysSinceLastUpload       uint8   `form:"expiryDaysSinceLastUpload"       binding:"required"`
	ExpiryTotalUploads              uint8   `form:"expiryTotalUploads"              binding:"required"`
	FileExpiryType                  string  `form:"fileExpiryType"                  binding:"required"`
	FileExpiryTotalDays             uint8   `form:"fileExpiryTotalDays"             binding:"required"`
	FileExpiryDaysSinceLastDownload uint8   `form:"fileExpiryDaysSinceLastDownload" binding:"required"`
	FileExpiryTotalDownloads        uint8   `form:"fileExpiryTotalDownloads"        binding:"required"`
	EmailOnUpload                   *string `form:"emailOnUpload"`
	CreatorLang                     string  `form:"creatorLang"                     binding:"required"`
	ReceiverLang                    string  `form:"receiverLang"                    binding:"required"`
}

func createGrantFactory(configValue config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := middleware.GetCurrentUser(c)
		var form grantForm
		tx, err := config.DBClient.Tx(c.Request.Context())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		if err := c.ShouldBind(&form); err == nil {
			salt, err := util.GenerateSalt()
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			hashedPassword := util.HashPassword(form.Password, salt)
			grantBuilder := tx.Grant.Create().
				SetID(uuid.New()).
				SetExpiryType(form.ExpiryType).
				SetExpiryDaysSinceLastUpload(form.ExpiryDaysSinceLastUpload).
				SetExpiryTotalDays(form.ExpiryTotalDays).
				SetExpiryTotalUploads(form.ExpiryTotalUploads).
				SetFileExpiryType(form.ExpiryType).
				SetFileExpiryDaysSinceLastDownload(form.FileExpiryDaysSinceLastDownload).
				SetFileExpiryTotalDays(form.FileExpiryTotalDays).
				SetFileExpiryTotalDownloads(form.FileExpiryTotalDownloads).
				SetHashedPassword(hashedPassword).
				SetSalt(hex.EncodeToString(salt)).
				SetOwner(currentUser).
				SetCreatorLang(form.CreatorLang)

			if form.Comment != nil {
				grantBuilder = grantBuilder.SetComment(*form.Comment)
			}

			if form.EmailOnUpload != nil {
				grantBuilder = grantBuilder.SetEmailOnUpload(*form.EmailOnUpload)
			}

			grantValue, err := grantBuilder.Save(c.Request.Context())
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
			}
			tx.User.UpdateOne(currentUser).AddSubmittedGrants(1).SaveX(c.Request.Context())
			grantValue = tx.Grant.Query().
				WithOwner().
				WithFiles().
				Where(grant.ID(grantValue.ID)).
				OnlyX(c.Request.Context())
			c.JSON(
				http.StatusCreated,
				apiTypes.ToPublicGrant(configValue, grantValue, []*ent.File{}),
			)
			if form.Email != nil {
				var toBeEmailedPassword *string = nil
				if form.EmailPassword {
					toBeEmailedPassword = &form.Password
				}
				mail.SendGrantSharedNotification(
					c,
					configValue,
					*form.Email,
					form.ReceiverLang,
					grantValue,
					toBeEmailedPassword,
				)
			}

			tx.Commit()
		} else {
			c.AbortWithError(422, err)
		}

	}
}

func fetchGrantsFactory(configValue config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := middleware.GetCurrentUser(c)
		query := config.DBClient.Grant.Query().WithFiles().WithOwner()

		if !currentUser.IsAdmin {
			query = query.Where(grant.HasOwnerWith(user.ID(currentUser.ID)))
		}

		grants, err := query.All(c.Request.Context())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		publicGrants := make([]apiTypes.PublicGrant, 0, len(grants))
		for _, grantValue := range grants {
			if !util.GrantExpired(configValue, grantValue) {
				publicGrants = append(
					publicGrants,
					apiTypes.ToPublicGrant(configValue, grantValue, grantValue.Edges.Files),
				)
			}
		}
		c.JSON(http.StatusOK, publicGrants)
	}
}

func setupGrantGroup(r *gin.RouterGroup, configValue config.Config) {
	r.POST("", createGrantFactory(configValue))
	r.GET("", fetchGrantsFactory(configValue))
}
