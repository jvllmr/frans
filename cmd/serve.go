package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/routes"
	"github.com/spf13/cobra"
)

func startGin(configValue config.Config, db *ent.Client) {

	if configValue.DevMode {
		slog.Info("frans was started in development mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r, err := routes.SetupRootRouter(configValue, db)
	if err != nil {
		slog.Error("Setup failed", "err", err)
		os.Exit(1)
	}
	serveString := fmt.Sprintf("%s:%d", configValue.Host, configValue.Port)
	slog.Info(fmt.Sprintf("Serving on %s", serveString))
	r.Run(serveString)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start only the server",
	Run: func(cmd *cobra.Command, args []string) {
		configValue, db := getConfigAndDBClient()
		defer db.Close()
		startGin(configValue, db)
	},
}
