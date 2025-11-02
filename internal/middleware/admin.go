package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/otel"
	"github.com/jvllmr/frans/internal/util"
)

func AdminRequired(c *gin.Context) {
	_, span := otel.NewSpan(c.Request.Context(), "adminRequired")
	defer span.End()
	currentUser := GetCurrentUser(c)
	if !currentUser.IsAdmin {
		util.GinAbortWithError(
			c,
			http.StatusForbidden,
			fmt.Errorf("user %s is not an administrator", currentUser.Username),
		)
		return
	}
}
