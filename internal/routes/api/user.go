package apiRoutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/middleware"
	apiTypes "github.com/jvllmr/frans/internal/routes/api/types"
)

func fetchMe(ctx *gin.Context) {
	currentUser := middleware.GetCurrentUser(ctx)
	activeTickets := currentUser.QueryTickets().CountX(ctx.Request.Context())

	ctx.JSON(http.StatusOK, apiTypes.ToAdminViewUser(currentUser, activeTickets, 0))
}

func fetchUsers(ctx *gin.Context) {
	publicUsers := make([]apiTypes.AdminViewUser, 0)
	users := config.DBClient.User.Query().AllX(ctx.Request.Context())
	for _, userValue := range users {
		activeTickets := userValue.QueryTickets().CountX(ctx.Request.Context())
		publicUsers = append(publicUsers, apiTypes.ToAdminViewUser(userValue, activeTickets, 0))
	}
	ctx.JSON(http.StatusOK, publicUsers)
}

func setupUserGroup(r *gin.RouterGroup, configValue config.Config) {
	r.GET("/me", fetchMe)
	r.GET("", middleware.AdminRequired(configValue), fetchUsers)
}
