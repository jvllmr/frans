package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/oidc"
)

func AdminRequired(p *oidc.FransOidcProvider) gin.HandlerFunc {
	auth := Auth(p, false)

	return func(c *gin.Context) {
		auth(c)
		currentUser := GetCurrentUser(c)
		if !currentUser.IsAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
