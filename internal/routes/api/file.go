package apiRoutes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
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
	db          *ent.Client
	fileService services.FileService
}

func (fc *fileController) fetchReceivedFilesHandler(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)

	files := fc.db.File.Query().
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
	fileValue, err := fc.db.File.Query().
		WithData().
		Where(file.ID(uuid.MustParse(requestedFile.ID))).
		Only(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	currentUser := middleware.GetCurrentUser(c)

	if !util.UserHasFileAccess(c.Request.Context(), currentUser, fileValue) {
		c.AbortWithStatus(http.StatusForbidden)
	}

	addDownload := c.Query("addDownload")
	if len(addDownload) > 0 {
		fc.db.File.UpdateOne(fileValue).
			AddTimesDownloaded(1).
			SetLastDownload(time.Now()).
			ExecX(c.Request.Context())
	}

	filePath := fc.fileService.FilesFilePath(fileValue.Edges.Data.ID)
	c.FileAttachment(filePath, fileValue.Name)
}

func setupFileGroup(r *gin.RouterGroup, configValue config.Config, db *ent.Client) {
	controller := fileController{
		config:      configValue,
		db:          db,
		fileService: services.NewFileService(configValue, db),
	}
	r.GET("/received", controller.fetchReceivedFilesHandler)
	r.GET("/:fileId", controller.fetchFileHandler)
}
