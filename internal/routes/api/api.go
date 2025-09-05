package apiRoutes

import (
	"github.com/gin-gonic/gin"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/middleware"
	"github.com/jvllmr/frans/internal/oidc"

	shareRoutes "github.com/jvllmr/frans/internal/routes/api/share"
)

func SetupAPIRoutes(r *gin.RouterGroup, configValue config.Config, pkceCache *oidc.PKCECache) {
	apiGroup := r.Group("/api")
	setupAuthRoutes(apiGroup, configValue, pkceCache)

	v1Group := apiGroup.Group("/v1")
	userGroup := v1Group.Group("/user", middleware.Auth(configValue, nil))
	setupUserGroup(userGroup, configValue)
	ticketGroup := v1Group.Group("/ticket", middleware.Auth(configValue, nil))
	setupTicketGroup(ticketGroup, configValue)
	grantGroup := v1Group.Group("/grant", middleware.Auth(configValue, nil))
	setupGrantGroup(grantGroup, configValue)
	fileGroup := v1Group.Group("/file", middleware.Auth(configValue, nil))
	setupFileGroup(fileGroup, configValue)
	shareGroup := v1Group.Group("/share")
	shareRoutes.SetupShareRoutes(shareGroup, configValue)
}
