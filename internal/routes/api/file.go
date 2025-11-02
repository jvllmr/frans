package apiRoutes

import (
	"fmt"
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
	"github.com/jvllmr/frans/internal/otel"
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
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchReceivedFiles")
	defer span.End()
	currentUser := middleware.GetCurrentUser(c)

	files := fc.db.File.Query().WithData().
		Where(file.HasGrantsWith(grant.HasOwnerWith(user.ID(currentUser.ID)))).
		AllX(ctx)

	publicFiles := make([]services.PublicFile, len(files))
	for i, fileValue := range files {
		publicFiles[i] = fc.fileService.ToPublicFile(fileValue)
	}

	c.JSON(http.StatusOK, publicFiles)
}

func (fc *fileController) fetchFileHandler(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchFile")
	defer span.End()
	var requestedFile apiTypes.RequestedFileParam
	if err := c.ShouldBindUri(&requestedFile); err != nil {
		util.GinAbortWithError(c, http.StatusBadRequest, err)
		return
	}
	fileValue, err := fc.db.File.Query().
		WithData().
		Where(file.ID(uuid.MustParse(requestedFile.ID))).
		Only(ctx)
	if err != nil {
		util.GinAbortWithError(c, http.StatusNotFound, err)
		return
	}

	currentUser := middleware.GetCurrentUser(c)

	if !util.UserHasFileAccess(ctx, currentUser, fileValue) {
		util.GinAbortWithError(
			c,
			http.StatusForbidden,
			fmt.Errorf("user %s does not have access to file", currentUser.Username),
		)
		return
	}

	addDownload := c.Query("addDownload")
	if len(addDownload) > 0 {
		fc.db.File.UpdateOne(fileValue).
			AddTimesDownloaded(1).
			SetLastDownload(time.Now()).
			ExecX(ctx)
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
