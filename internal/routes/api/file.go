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
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/util"
)

type fileController struct {
	config      config.Config
	fileService services.FileService
}

func (fc *fileController) fetchReceivedFilesHandler(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)

	files := config.DBClient.File.Query().
		Where(file.HasGrantsWith(grant.HasOwnerWith(user.ID(currentUser.ID)))).
		AllX(c.Request.Context())

	publicFiles := make([]services.PublicFile, len(files))
	for i, fileValue := range files {
		publicFiles[i] = fc.fileService.ToPublicFile(fileValue)
	}

	c.JSON(http.StatusOK, publicFiles)
}

func (fc *fileController) fetchFileHandler(c *gin.Context) {

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
	filePath := fc.fileService.FilesFilePath(fileValue.Sha512)
	c.FileAttachment(filePath, fileValue.Name)

}

func setupFileGroup(r *gin.RouterGroup, configValue config.Config) {
	controller := fileController{
		config:      configValue,
		fileService: services.NewFileService(configValue),
	}
	r.GET("/received", controller.fetchReceivedFilesHandler)
	r.GET("/:fileId", controller.fetchFileHandler)
}
