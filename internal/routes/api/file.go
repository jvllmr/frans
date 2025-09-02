package apiRoutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent/file"
	"github.com/jvllmr/frans/internal/ent/grant"
	"github.com/jvllmr/frans/internal/ent/user"
	"github.com/jvllmr/frans/internal/middleware"
	apiTypes "github.com/jvllmr/frans/internal/routes/api/types"
	"github.com/jvllmr/frans/internal/util"
)

func fetchReceivedFilesRouteFactory(configValue config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := middleware.GetCurrentUser(c)

		files := config.DBClient.File.Query().
			Where(file.HasGrantsWith(grant.HasOwnerWith(user.ID(currentUser.ID)))).
			AllX(c.Request.Context())

		publicFiles := make([]apiTypes.PublicFile, len(files))
		for i, fileValue := range files {
			publicFiles[i] = apiTypes.ToPublicFile(configValue, fileValue)
		}

		c.JSON(http.StatusOK, publicFiles)
	}
}

func fetchFileRouteFactory(configValue config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

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

		currentUser := middleware.GetCurrentUser(c)

		if !util.UserHasFileAccess(c.Request.Context(), currentUser, fileValue) {
			c.AbortWithStatus(http.StatusForbidden)
		}
		filePath := util.FilesFilePath(configValue, fileValue.Sha512)
		c.FileAttachment(filePath, fileValue.Name)

	}
}

func setupFileGroup(r *gin.RouterGroup, configValue config.Config) {
	r.GET("/received", fetchReceivedFilesRouteFactory(configValue))
	r.GET("/:fileId", fetchFileRouteFactory(configValue))
}
