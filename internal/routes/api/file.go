package apiRoutes

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/file"
	"github.com/jvllmr/frans/internal/ent/user"
	"github.com/jvllmr/frans/internal/mail"
	"github.com/jvllmr/frans/internal/middleware"
	"github.com/jvllmr/frans/internal/otel"
	apiTypes "github.com/jvllmr/frans/internal/routes/api/types"
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/util"
)

type fileController struct {
	cfg         config.Config
	db          *ent.Client
	fileService services.FileService
	mailer      mail.Mailer
}

func (fc *fileController) fetchReceivedFilesHandler(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchReceivedFiles")
	defer span.End()
	currentUser := middleware.GetCurrentUser(c)

	filesQuery := fc.db.File.Query().WithData().WithOwner().Order(file.ByCreatedAt(sql.OrderDesc()))

	if !currentUser.IsAdmin {
		filesQuery = filesQuery.Where(
			file.HasOwnerWith(user.ID(currentUser.ID)),
			sql.NotPredicates(file.HasTicket()),
		)
	}

	files := filesQuery.AllX(ctx)

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
		util.GinAbortWithError(ctx, c, http.StatusBadRequest, err)
		return
	}
	fileValue, err := fc.db.File.Query().
		WithData().WithOwner().
		Where(file.ID(uuid.MustParse(requestedFile.ID))).
		Only(ctx)
	if err != nil {
		util.GinAbortWithError(ctx, c, http.StatusNotFound, err)
		return
	}

	currentUser := middleware.GetCurrentUser(c)

	if !util.UserHasFileAccess(ctx, currentUser, fileValue) {
		util.GinAbortWithError(ctx,
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

func (fc *fileController) deleteFileHandler(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "deleteFileManual")
	defer span.End()
	var requestedFile apiTypes.RequestedFileParam
	if err := c.ShouldBindUri(&requestedFile); err != nil {
		util.GinAbortWithError(ctx, c, http.StatusBadRequest, err)
		return
	}
	f, err := fc.db.File.Query().
		Where(file.ID(uuid.MustParse(requestedFile.ID))).
		WithOwner().
		WithData().
		WithGrant().
		WithTicket().
		Only(ctx)
	if err != nil {
		util.GinAbortWithError(ctx, c, http.StatusNotFound, err)
		return
	}
	currentUser := middleware.GetCurrentUser(c)
	isUserOwner := f.Edges.Owner.ID == currentUser.ID
	if !currentUser.IsAdmin && !isUserOwner {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err := fc.fileService.DeleteFile(ctx, f); err != nil {
		util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
		return
	}

	if !isUserOwner {
		if err := fc.mailer.SendFileDeletionNotification(f, fc.cfg.GetBaseURL(c.Request)); err != nil {
			util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
			return
		}
	}
	slog.InfoContext(
		ctx,
		"Manual file deletion",
		"username",
		currentUser.Username,
		"owner",
		f.Edges.Owner.Username,
		"fileId",
		f.ID.String(),
	)
	c.Status(http.StatusOK)
}

func setupFileGroup(r *gin.RouterGroup, cfg config.Config, db *ent.Client) {
	controller := fileController{
		cfg:         cfg,
		db:          db,
		fileService: services.NewFileService(cfg, db),
		mailer:      mail.NewMailer(cfg),
	}
	r.GET("/received", controller.fetchReceivedFilesHandler)
	r.GET("/:fileId", controller.fetchFileHandler)
	r.DELETE("/:fileId", controller.deleteFileHandler)
}
