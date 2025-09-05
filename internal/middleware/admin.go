package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
)

func AdminRequired(configValue config.Config) gin.HandlerFunc {
	auth := Auth(configValue, nil)

	return func(c *gin.Context) {
		auth(c)
		currentUser := GetCurrentUser(c)
		if !currentUser.IsAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
