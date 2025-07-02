package apiRoutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent"
)

type PublicUser struct {
	ID       string `json:"id"`
	FullName string `json:"name"`
}

func ToPublicUser(user *ent.User) PublicUser {
	return PublicUser{
		ID:       user.ID.String(),
		FullName: user.FullName,
	}
}

func fetchMe(ctx *gin.Context) {
	user := ctx.MustGet(config.UserGinContext).(*ent.User)

	ctx.JSON(http.StatusOK, ToPublicUser(user))
}

func setupUserGroup(r *gin.RouterGroup) {
	r.GET("/me", fetchMe)
}
