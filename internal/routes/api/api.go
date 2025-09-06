package apiRoutes

import (
	"github.com/gin-gonic/gin"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/middleware"
	"github.com/jvllmr/frans/internal/oidc"

	shareRoutes "github.com/jvllmr/frans/internal/routes/api/share"
)

func SetupAPIRoutes(
	r *gin.RouterGroup,
	configValue config.Config,
	oidcProvider *oidc.FransOidcProvider,
) {
	apiGroup := r.Group("/api")
	setupAuthRoutes(apiGroup, configValue, oidcProvider)

	auth := middleware.Auth(oidcProvider, false)

	v1Group := apiGroup.Group("/v1")

	userGroup := v1Group.Group("/user", auth)
	setupUserGroup(userGroup, oidcProvider)

	ticketGroup := v1Group.Group("/ticket", auth)
	setupTicketGroup(ticketGroup, configValue)

	grantGroup := v1Group.Group("/grant", auth)
	setupGrantGroup(grantGroup, configValue)

	fileGroup := v1Group.Group("/file", auth)
	setupFileGroup(fileGroup, configValue)

	shareGroup := v1Group.Group("/share")
	shareRoutes.SetupShareRoutes(shareGroup, configValue)
}
