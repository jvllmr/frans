package apiRoutes

import (
	"github.com/gin-gonic/gin"

	"github.com/jvllmr/frans/pkg/config"
	routesUtil "github.com/jvllmr/frans/pkg/routes"
)

func SetupAPIRoutes(r *gin.RouterGroup, configValue config.Config) {
	apiGroup := r.Group("/api")
	setupAuthRoutes(apiGroup, configValue)

	v1Group := apiGroup.Group("/v1", routesUtil.AuthMiddleware(configValue, false))
	userGroup := v1Group.Group("/user")
	setupUserGroup(userGroup)
}
