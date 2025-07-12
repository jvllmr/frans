package apiRoutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent"
	apiTypes "github.com/jvllmr/frans/pkg/routes/api/types"
)

func fetchMe(ctx *gin.Context) {
	user := ctx.MustGet(config.UserGinContext).(*ent.User)
	ctx.JSON(http.StatusOK, apiTypes.ToPublicUser(user))
}

func setupUserGroup(r *gin.RouterGroup) {
	r.GET("/me", fetchMe)
}
