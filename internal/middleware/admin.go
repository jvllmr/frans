package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/otel"
)

func AdminRequired(c *gin.Context) {
	_, span := otel.NewSpan(c.Request.Context(), "adminRequired")
	defer span.End()
	currentUser := GetCurrentUser(c)
	if !currentUser.IsAdmin {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}
