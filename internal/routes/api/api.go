package apiRoutes

import (
	"github.com/gin-gonic/gin"

	"codeberg.org/jvllmr/frans/internal/config"
	"codeberg.org/jvllmr/frans/internal/ent"
	"codeberg.org/jvllmr/frans/internal/middleware"
	"codeberg.org/jvllmr/frans/internal/oidc"

	shareRoutes "codeberg.org/jvllmr/frans/internal/routes/api/share"
)

func SetupAPIRoutes(
	r *gin.RouterGroup,
	configValue config.Config,
	db *ent.Client,
	oidcProvider *oidc.FransOidcProvider,
) {
	apiGroup := r.Group("/api")
	setupAuthRoutes(apiGroup, configValue, oidcProvider)

	auth := middleware.Auth(oidcProvider, false)

	v1Group := apiGroup.Group("/v1")

	userGroup := v1Group.Group("/user", auth)
	setupUserGroup(userGroup, db)

	ticketGroup := v1Group.Group("/ticket", auth)
	setupTicketGroup(ticketGroup, configValue, db)

	grantGroup := v1Group.Group("/grant", auth)
	setupGrantGroup(grantGroup, configValue, db)

	fileGroup := v1Group.Group("/file", auth)
	setupFileGroup(fileGroup, configValue, db)

	shareGroup := v1Group.Group("/share")
	shareRoutes.SetupShareRoutes(shareGroup, configValue, db)
}
