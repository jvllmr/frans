package middleware

import (
	"fmt"
	"net/http"

	"codeberg.org/jvllmr/frans/internal/otel"
	"codeberg.org/jvllmr/frans/internal/util"
	"github.com/gin-gonic/gin"
)

func AdminRequired(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "adminRequired")
	defer span.End()
	currentUser := GetCurrentUser(c)
	if !currentUser.IsAdmin {
		util.GinAbortWithError(
			ctx,
			c,
			http.StatusForbidden,
			fmt.Errorf("user %s is not an administrator", currentUser.Username),
		)
		return
	}
}
