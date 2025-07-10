package apiRoutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent"
	"github.com/jvllmr/frans/pkg/ent/ticket"
	"github.com/jvllmr/frans/pkg/ent/user"
	"github.com/jvllmr/frans/pkg/util"
)

type PublicFile struct {
	Id              uuid.UUID `json:"id"`
	Sha512          string    `json:"sha512"`
	Size            uint64    `json:"size"`
	Name            string    `json:"name"`
	TimesDownloaded uint64    `json:"timesDownloaded"`
	LastDownloaded  *string   `json:"lastDownloaded"`
}

func ToPublicFile(file *ent.File) PublicFile {
	var lastDownloadedValue *string = nil
	if file.LastDownload != nil {
		formattedValue := file.LastDownload.UTC().Format(http.TimeFormat)
		lastDownloadedValue = &formattedValue
	}

	return PublicFile{
		Id:              file.ID,
		Sha512:          file.Sha512,
		Size:            file.Size,
		Name:            file.Name,
		TimesDownloaded: file.TimesDownloaded,
		LastDownloaded:  lastDownloadedValue,
	}

}

type RequestedFile struct {
	ID string `uri:"fileId" binding:"required,uuid"`
}

func fetchFileRouteFactory(configValue config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		var requestedFile RequestedFile
		if err := c.ShouldBindUri(&requestedFile); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		file, err := config.DBClient.File.Get(c.Request.Context(), uuid.MustParse(requestedFile.ID))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}

		currentUser := c.MustGet(config.UserGinContext).(*ent.User)

		if !currentUser.IsAdmin &&
			file.QueryTickets().
				Where(ticket.HasOwnerWith(user.ID(currentUser.ID))).
				CountX(c.Request.Context()) ==
				0 {
			c.AbortWithStatus(http.StatusForbidden)
		}
		filePath := util.GetFilesFilePath(configValue, file.Sha512)
		c.FileAttachment(filePath, file.Name)

	}
}

func setupFileGroup(r *gin.RouterGroup, configValue config.Config) {
	r.GET("/:fileId", fetchFileRouteFactory(configValue))
}
