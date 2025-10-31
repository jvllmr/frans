package apiRoutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/middleware"
	"github.com/jvllmr/frans/internal/otel"
	"github.com/jvllmr/frans/internal/services"
)

type userController struct {
	db *ent.Client
}

func (uc *userController) fetchMe(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchMe")
	defer span.End()
	currentUser := middleware.GetCurrentUser(c)
	activeTickets := currentUser.QueryTickets().CountX(ctx)
	activeGrants := currentUser.QueryGrants().CountX(ctx)

	c.JSON(http.StatusOK, services.ToAdminViewUser(currentUser, activeTickets, activeGrants))
}

func (uc *userController) fetchUsers(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchUsers")
	defer span.End()
	publicUsers := make([]services.AdminViewUser, 0)
	users := uc.db.User.Query().AllX(ctx)
	for _, userValue := range users {
		activeTickets := userValue.QueryTickets().CountX(ctx)
		activeGrants := userValue.QueryGrants().CountX(ctx)
		publicUsers = append(
			publicUsers,
			services.ToAdminViewUser(userValue, activeTickets, activeGrants),
		)
	}
	c.JSON(http.StatusOK, publicUsers)
}

func setupUserGroup(r *gin.RouterGroup, db *ent.Client) {
	controller := userController{db: db}

	r.GET("/me", controller.fetchMe)
	r.GET("", middleware.AdminRequired, controller.fetchUsers)
}
