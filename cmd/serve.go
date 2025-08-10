package cmd

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/routes"
	"github.com/spf13/cobra"
)

func startGin() {
	configValue := config.GetSafeConfig()
	config.InitOIDC(configValue)
	if configValue.DevMode {
		slog.Info("frans was started in development mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := routes.SetupRootRouter(configValue)
	serveString := fmt.Sprintf("%s:%d", configValue.Host, configValue.Port)
	slog.Info(fmt.Sprintf("Serving on %s", serveString))
	r.Run(serveString)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start only the server",
	Run:   func(cmd *cobra.Command, args []string) { startGin() },
}
