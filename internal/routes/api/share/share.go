package shareRoutes

import (
	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
)

func SetupShareRoutes(r *gin.RouterGroup, configValue config.Config) {
	ticketShareGroup := r.Group("/ticket")
	setupTicketShareRoutes(ticketShareGroup, configValue)
}
