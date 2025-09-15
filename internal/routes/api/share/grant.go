package shareRoutes

import (
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/grant"
	"github.com/jvllmr/frans/internal/mail"
	apiTypes "github.com/jvllmr/frans/internal/routes/api/types"
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/util"
)

type grantShareController struct {
	config       config.Config
	db           *ent.Client
	grantService services.GrantService
	fileService  services.FileService
	mailer       mail.Mailer
}

func (gsc *grantShareController) fetchGrant(c *gin.Context) {
	grantValue := c.MustGet(config.ShareGrantContext).(*ent.Grant)
	c.JSON(
		http.StatusOK,
		gsc.grantService.ToPublicGrant(gsc.fileService, grantValue, grantValue.Edges.Files),
	)
}

func (gsc *grantShareController) fetchGrantAccessToken(c *gin.Context) {
	tokenValueBytes := util.GenerateSalt()

	tokenValue := hex.EncodeToString(tokenValueBytes)
	token := gsc.db.ShareAccessToken.Create().
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
}

func (gsc *grantShareController) postGrantFiles(c *gin.Context) {
	multipartForm, _ := c.MultipartForm()
	files := multipartForm.File["files[]"]
	grantValue := c.MustGet(config.ShareGrantContext).(*ent.Grant)

	if len(files) > int(gsc.config.MaxFiles) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	gsc.fileService.EnsureFilesTmpPath()
	tx, err := gsc.db.Tx(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	dbFiles := make([]*ent.File, len(files))
	for i, fileHeader := range files {
		dbFile, err := gsc.fileService.CreateFile(
			c.Request.Context(),
			tx, fileHeader, grantValue.Edges.Owner,
			grantValue.ExpiryType,
			grantValue.FileExpiryDaysSinceLastDownload,
			grantValue.FileExpiryTotalDays,
			grantValue.FileExpiryTotalDownloads,
		)
		if err != nil {
			var errFileTooBig *services.ErrFileTooBig
			if errors.As(err, &errFileTooBig) {
				c.AbortWithError(http.StatusBadRequest, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}

		tx.Grant.UpdateOne(grantValue).
			AddFiles(dbFile).
			SaveX(c.Request.Context())
		dbFiles[i] = dbFile
	}
	tx.Grant.UpdateOne(grantValue).
		SetLastUpload(time.Now()).
		AddTimesUploaded(1).
		SaveX(c.Request.Context())

	err = util.RefreshUserTotalDataSize(c.Request.Context(), grantValue.Edges.Owner, tx)
	if err != nil {
		slog.Error(
			"Could not refresh total data size of user",
			"err",
			err,
			"user",
			grantValue.Edges.Owner,
		)
	}
	if grantValue.EmailOnUpload != nil && grantValue.TimesUploaded == 1 {
		gsc.mailer.SendFileUploadNotification(
			c,
			*grantValue.EmailOnUpload,
			grantValue,
			dbFiles,
		)
	}
	tx.Commit()
	grantValue = gsc.db.Grant.Query().
		Where(grant.ID(grantValue.ID)).
		WithFiles(func(fq *ent.FileQuery) { fq.WithData() }).
		WithOwner().
		OnlyX(c.Request.Context())

	c.JSON(
		http.StatusOK,
		gsc.grantService.ToPublicGrant(gsc.fileService, grantValue, grantValue.Edges.Files),
	)
}

func setupGrantShareRoutes(r *gin.RouterGroup, configValue config.Config, db *ent.Client) {
	getGrantMiddleware := func(c *gin.Context) {
		grantId := c.Param("grantId")
		username, password, ok := c.Request.BasicAuth()

		if !ok {
			tokenCookie, err := c.Cookie(config.ShareAccessTokenCookieName)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			token, err := db.ShareAccessToken.Get(c.Request.Context(), tokenCookie)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			} else if token.Expiry.Before(time.Now()) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else if username != grantId {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		uuidValue, err := uuid.Parse(grantId)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		grantValue, err := db.Grant.Query().
			Where(grant.ID(uuidValue)).
			WithOwner().
			WithFiles(func(fq *ent.FileQuery) { fq.WithData() }).
			Only(c.Request.Context())

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if ok && !util.VerifyPassword(password, grantValue.HashedPassword, grantValue.Salt) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(config.ShareGrantContext, grantValue)
	}

	singleGrantShareGroup := r.Group("/:grantId", getGrantMiddleware)

	controller := grantShareController{
		config:       configValue,
		db:           db,
		grantService: services.NewGrantService(configValue),
		fileService:  services.NewFileService(configValue, db),
		mailer:       mail.NewMailer(configValue),
	}

	singleGrantShareGroup.GET("", controller.fetchGrant)

	singleGrantShareGroup.GET("/token", controller.fetchGrantAccessToken)

	singleGrantShareGroup.POST("", controller.postGrantFiles)
}
