package shareRoutes

import (
	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
)

func SetupShareRoutes(r *gin.RouterGroup, configValue config.Config, db *ent.Client) {
	ticketShareGroup := r.Group("/ticket")
	grantShareGroup := r.Group("/grant")
	setupTicketShareRoutes(ticketShareGroup, configValue, db)
	setupGrantShareRoutes(grantShareGroup, configValue, db)
}
