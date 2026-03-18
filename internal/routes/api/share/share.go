package shareRoutes

import (
	"codeberg.org/jvllmr/frans/internal/config"
	"codeberg.org/jvllmr/frans/internal/ent"
	"github.com/gin-gonic/gin"
)

func SetupShareRoutes(r *gin.RouterGroup, configValue config.Config, db *ent.Client) {
	ticketShareGroup := r.Group("/ticket")
	grantShareGroup := r.Group("/grant")
	setupTicketShareRoutes(ticketShareGroup, configValue, db)
	setupGrantShareRoutes(grantShareGroup, configValue, db)
}
