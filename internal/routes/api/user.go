package apiRoutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/middleware"
	"github.com/jvllmr/frans/internal/oidc"
	"github.com/jvllmr/frans/internal/services"
)

type userController struct {
	db *ent.Client
}

func (uc *userController) fetchMe(ctx *gin.Context) {
	currentUser := middleware.GetCurrentUser(ctx)
	activeTickets := currentUser.QueryTickets().CountX(ctx.Request.Context())

	ctx.JSON(http.StatusOK, services.ToAdminViewUser(currentUser, activeTickets, 0))
}

func (uc *userController) fetchUsers(ctx *gin.Context) {
	publicUsers := make([]services.AdminViewUser, 0)
	users := uc.db.User.Query().AllX(ctx.Request.Context())
	for _, userValue := range users {
		activeTickets := userValue.QueryTickets().CountX(ctx.Request.Context())
		publicUsers = append(publicUsers, services.ToAdminViewUser(userValue, activeTickets, 0))
	}
	ctx.JSON(http.StatusOK, publicUsers)
}

func setupUserGroup(r *gin.RouterGroup, db *ent.Client, oidcProvider *oidc.FransOidcProvider) {
	controller := userController{db: db}

	r.GET("/me", controller.fetchMe)
	r.GET("", middleware.AdminRequired(oidcProvider), controller.fetchUsers)
}
