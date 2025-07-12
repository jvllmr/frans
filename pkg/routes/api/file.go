package apiRoutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent"
	"github.com/jvllmr/frans/pkg/ent/ticket"
	"github.com/jvllmr/frans/pkg/ent/user"
	apiTypes "github.com/jvllmr/frans/pkg/routes/api/types"
	"github.com/jvllmr/frans/pkg/util"
)

func fetchFileRouteFactory(configValue config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		var requestedFile apiTypes.RequestedFileParam
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
