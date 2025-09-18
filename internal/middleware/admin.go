package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRequired(c *gin.Context) {
	currentUser := GetCurrentUser(c)
	if !currentUser.IsAdmin {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}
