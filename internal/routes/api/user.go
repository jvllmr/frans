package apiRoutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/middleware"
	"github.com/jvllmr/frans/internal/oidc"
	"github.com/jvllmr/frans/internal/services"
)

func fetchMe(ctx *gin.Context) {
	currentUser := middleware.GetCurrentUser(ctx)
	activeTickets := currentUser.QueryTickets().CountX(ctx.Request.Context())

	ctx.JSON(http.StatusOK, services.ToAdminViewUser(currentUser, activeTickets, 0))
}

func fetchUsers(ctx *gin.Context) {
	publicUsers := make([]services.AdminViewUser, 0)
	users := config.DBClient.User.Query().AllX(ctx.Request.Context())
	for _, userValue := range users {
		activeTickets := userValue.QueryTickets().CountX(ctx.Request.Context())
		publicUsers = append(publicUsers, services.ToAdminViewUser(userValue, activeTickets, 0))
	}
	ctx.JSON(http.StatusOK, publicUsers)
}

func setupUserGroup(r *gin.RouterGroup, oidcProvider *oidc.FransOidcProvider) {
	r.GET("/me", fetchMe)
	r.GET("", middleware.AdminRequired(oidcProvider), fetchUsers)
}
