package shareRoutes

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/grant"
	"github.com/jvllmr/frans/internal/mail"
	apiTypes "github.com/jvllmr/frans/internal/routes/api/types"
	"github.com/jvllmr/frans/internal/util"
)

func getGrantMiddleware(c *gin.Context) {
	grantId := c.Param("grantId")
	username, password, ok := c.Request.BasicAuth()

	if !ok {
		tokenCookie, err := c.Cookie(config.ShareAccessTokenCookieName)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token, err := config.DBClient.ShareAccessToken.Get(c.Request.Context(), tokenCookie)
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
	grantValue, err := config.DBClient.Grant.Query().
		Where(grant.ID(uuidValue)).
		WithOwner().
		WithFiles().
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

func setupGrantShareRoutes(r *gin.RouterGroup, configValue config.Config) {
	singleGrantShareGroup := r.Group("/:grantId", getGrantMiddleware)

	singleGrantShareGroup.GET("", func(c *gin.Context) {
		grantValue := c.MustGet(config.ShareGrantContext).(*ent.Grant)
		c.JSON(
			http.StatusOK,
			apiTypes.ToPublicGrant(configValue, grantValue, grantValue.Edges.Files),
		)
	})

	singleGrantShareGroup.GET("/token", func(c *gin.Context) {
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

	singleGrantShareGroup.POST("", func(c *gin.Context) {
		multipartForm, _ := c.MultipartForm()
		files := multipartForm.File["files[]"]
		grantValue := c.MustGet(config.ShareGrantContext).(*ent.Grant)

		if len(files) > int(configValue.MaxFiles) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		util.EnsureFilesTmpPath(configValue)
		tx, err := config.DBClient.Tx(c.Request.Context())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		dbFiles := make([]*ent.File, len(files))
		for i, fileHeader := range files {
			if fileHeader.Size > configValue.MaxSizes {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			incomingFileHandle, _ := fileHeader.Open()
			hasher := sha512.New()
			tmpFilePath := util.FilesTmpFilePath(configValue)
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
				SetExpiryType(grantValue.FileExpiryType).
				SetExpiryDaysSinceLastDownload(grantValue.FileExpiryDaysSinceLastDownload).
				SetExpiryTotalDays(grantValue.FileExpiryTotalDays).
				SetExpiryTotalDownloads(grantValue.FileExpiryTotalDownloads).
				Save(c.Request.Context())
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			targetFilePath := util.FilesFilePath(configValue, sha512sum)
			if _, err = os.Stat(targetFilePath); err != nil {
				os.Rename(tmpFilePath, targetFilePath)
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

		err = util.RefreshUserTotalDataSize(c.Request.Context(), grantValue.Edges.Owner)
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
			mail.SendFileUploadNotification(
				c,
				configValue,
				*grantValue.EmailOnUpload,
				grantValue,
				dbFiles,
			)
		}
		tx.Commit()
		grantValue = config.DBClient.Grant.Query().
			Where(grant.ID(grantValue.ID)).
			WithFiles().
			WithOwner().
			OnlyX(c.Request.Context())

		c.JSON(
			http.StatusOK,
			apiTypes.ToPublicGrant(configValue, grantValue, grantValue.Edges.Files),
		)
	})

}
