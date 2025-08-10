package apiRoutes

import (
	"github.com/gin-gonic/gin"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/middleware"

	shareRoutes "github.com/jvllmr/frans/internal/routes/api/share"
)

func SetupAPIRoutes(r *gin.RouterGroup, configValue config.Config) {
	apiGroup := r.Group("/api")
	setupAuthRoutes(apiGroup, configValue)

	v1Group := apiGroup.Group("/v1")
	userGroup := v1Group.Group("/user", middleware.Auth(configValue, false))
	setupUserGroup(userGroup)
	ticketGroup := v1Group.Group("/ticket", middleware.Auth(configValue, false))
	setupTicketGroup(ticketGroup, configValue)
	fileGroup := v1Group.Group("/file", middleware.Auth(configValue, false))
	setupFileGroup(fileGroup, configValue)
	shareGroup := v1Group.Group("/share")
	shareRoutes.SetupShareRoutes(shareGroup, configValue)
}
