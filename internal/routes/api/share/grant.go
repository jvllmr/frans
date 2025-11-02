package shareRoutes

import (
	"encoding/hex"
	"errors"
	"fmt"
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
	"github.com/jvllmr/frans/internal/otel"
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
	_, span := otel.NewSpan(c.Request.Context(), "fetchGrantShare")
	defer span.End()
	grantValue := c.MustGet(config.ShareGrantContext).(*ent.Grant)
	c.JSON(
		http.StatusOK,
		gsc.grantService.ToPublicGrant(gsc.fileService, grantValue, grantValue.Edges.Files),
	)
}

func (gsc *grantShareController) fetchGrantAccessToken(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchGrantShareAccessToken")
	defer span.End()
	tokenValueBytes := util.GenerateSalt()

	tokenValue := hex.EncodeToString(tokenValueBytes)
	token := gsc.db.ShareAccessToken.Create().
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

func (gsc *grantShareController) postGrantFiles(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "postGrantShareFiles")
	defer span.End()
	multipartForm, _ := c.MultipartForm()
	files := multipartForm.File["files[]"]
	grantValue := c.MustGet(config.ShareGrantContext).(*ent.Grant)

	if len(files) > int(gsc.config.MaxFiles) {
		util.GinAbortWithError(
			c,
			http.StatusBadRequest,
			fmt.Errorf(
				"maximum of %d files allowed per upload. %d uploaded",
				gsc.config.MaxFiles,
				len(files),
			),
		)
		return
	}
	gsc.fileService.EnsureFilesTmpPath()
	tx, err := gsc.db.Tx(ctx)
	if err != nil {
		util.GinAbortWithError(c, http.StatusInternalServerError, err)
	}

	dbFiles := make([]*ent.File, len(files))
	for i, fileHeader := range files {
		dbFile, err := gsc.fileService.CreateFile(
			ctx,
			tx, fileHeader, grantValue.Edges.Owner,
			grantValue.ExpiryType,
			grantValue.FileExpiryDaysSinceLastDownload,
			grantValue.FileExpiryTotalDays,
			grantValue.FileExpiryTotalDownloads,
		)
		if err != nil {
			var errFileTooBig *services.ErrFileTooBig
			if errors.As(err, &errFileTooBig) {
				util.GinAbortWithError(c, http.StatusBadRequest, err)
			} else {
				util.GinAbortWithError(c, http.StatusInternalServerError, err)
			}
		}

		tx.Grant.UpdateOne(grantValue).
			AddFiles(dbFile).
			SaveX(ctx)
		dbFiles[i] = dbFile
	}
	tx.Grant.UpdateOne(grantValue).
		SetLastUpload(time.Now()).
		AddTimesUploaded(1).
		SaveX(ctx)

	err = util.RefreshUserTotalDataSize(ctx, grantValue.Edges.Owner, tx)
	if err != nil {
		slog.Error(
			"Could not refresh total data size of user",
			"err",
			err,
			"user",
			grantValue.Edges.Owner,
		)
	}
	// we have an outdated reference to the grant; therefore we check if TimesUploaded is 0
	if grantValue.EmailOnUpload != nil && grantValue.TimesUploaded == 0 {
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
		OnlyX(ctx)

	c.JSON(
		http.StatusOK,
		gsc.grantService.ToPublicGrant(gsc.fileService, grantValue, grantValue.Edges.Files),
	)
}

func setupGrantShareRoutes(r *gin.RouterGroup, configValue config.Config, db *ent.Client) {
	getGrantMiddleware := func(c *gin.Context) {
		ctx, span := otel.NewSpan(c.Request.Context(), "checkGrantShareAuth")
		defer span.End()
		grantId := c.Param("grantId")
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
				return
			} else if token.Expiry.Before(time.Now()) {
				util.GinAbortWithError(c, http.StatusUnauthorized, fmt.Errorf("token expired"))
				return
			}
		} else if username != grantId {
			util.GinAbortWithError(c, http.StatusUnauthorized, fmt.Errorf("token does not match share"))
			return
		}
		uuidValue, err := uuid.Parse(grantId)
		if err != nil {
			util.GinAbortWithError(c, http.StatusBadRequest, err)
			return
		}
		grantValue, err := db.Grant.Query().
			Where(grant.ID(uuidValue)).
			WithOwner().
			WithFiles(func(fq *ent.FileQuery) { fq.WithData() }).
			Only(ctx)

		if err != nil {
			util.GinAbortWithError(c, http.StatusUnauthorized, err)
			return
		}

		if ok && !util.VerifyPassword(password, grantValue.HashedPassword, grantValue.Salt) {
			util.GinAbortWithError(c, http.StatusUnauthorized, fmt.Errorf("password incorrect"))
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
