package apiRoutes

import (
	"github.com/gin-gonic/gin"

	"github.com/jvllmr/frans/pkg/config"
	routesUtil "github.com/jvllmr/frans/pkg/routes"
	shareRoutes "github.com/jvllmr/frans/pkg/routes/api/share"
)

func SetupAPIRoutes(r *gin.RouterGroup, configValue config.Config) {
	apiGroup := r.Group("/api")
	setupAuthRoutes(apiGroup, configValue)

	v1Group := apiGroup.Group("/v1")
	userGroup := v1Group.Group("/user", routesUtil.AuthMiddleware(configValue, false))
	setupUserGroup(userGroup)
	ticketGroup := v1Group.Group("/ticket", routesUtil.AuthMiddleware(configValue, false))
	setupTicketGroup(ticketGroup, configValue)
	fileGroup := v1Group.Group("/file", routesUtil.AuthMiddleware(configValue, false))
	setupFileGroup(fileGroup, configValue)
	shareGroup := v1Group.Group("/share")
	shareRoutes.SetupShareRoutes(shareGroup, configValue)
}
