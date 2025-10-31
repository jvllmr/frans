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
	"github.com/jvllmr/frans/internal/otel"
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/util"
)

type grantController struct {
	config       config.Config
	db           *ent.Client
	grantService services.GrantService
	fileService  services.FileService
	mailer       mail.Mailer
}

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

func (gc *grantController) createGrantHandler(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "createGrant")
	defer span.End()
	currentUser := middleware.GetCurrentUser(c)
	var form grantForm
	tx, err := gc.db.Tx(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	if err := c.ShouldBind(&form); err == nil {
		salt := util.GenerateSalt()

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

		grantValue, err := grantBuilder.Save(ctx)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		tx.User.UpdateOne(currentUser).AddSubmittedGrants(1).SaveX(ctx)
		grantValue = tx.Grant.Query().
			WithOwner().
			WithFiles(func(fq *ent.FileQuery) {
				fq.WithData()
			}).
			Where(grant.ID(grantValue.ID)).
			OnlyX(ctx)
		c.JSON(
			http.StatusCreated,
			gc.grantService.ToPublicGrant(gc.fileService, grantValue, make([]*ent.File, 0)),
		)
		if form.Email != nil {
			var toBeEmailedPassword *string = nil
			if form.EmailPassword {
				toBeEmailedPassword = &form.Password
			}
			gc.mailer.SendGrantSharedNotification(
				c,
				gc.grantService,
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

func (gc *grantController) fetchGrantsHandler(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchGrants")
	defer span.End()
	currentUser := middleware.GetCurrentUser(c)
	query := gc.db.Grant.Query().WithFiles(func(fq *ent.FileQuery) { fq.WithData() }).WithOwner()

	if !currentUser.IsAdmin {
		query = query.Where(grant.HasOwnerWith(user.ID(currentUser.ID)))
	}

	grants, err := query.All(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	publicGrants := make([]services.PublicGrant, 0, len(grants))
	for _, grantValue := range grants {
		publicGrants = append(
			publicGrants,
			gc.grantService.ToPublicGrant(gc.fileService, grantValue, grantValue.Edges.Files),
		)
	}
	c.JSON(http.StatusOK, publicGrants)
}

func setupGrantGroup(r *gin.RouterGroup, configValue config.Config, db *ent.Client) {
	controller := grantController{
		config:       configValue,
		db:           db,
		grantService: services.NewGrantService(configValue),
		fileService:  services.NewFileService(configValue, db),
		mailer:       mail.NewMailer(configValue),
	}
	r.POST("", controller.createGrantHandler)
	r.GET("", controller.fetchGrantsHandler)
}
