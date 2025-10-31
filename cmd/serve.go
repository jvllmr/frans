package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/otel"
	"github.com/jvllmr/frans/internal/routes"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func startGin(ctx context.Context, configValue config.Config, db *ent.Client) {

	if configValue.DevMode {
		slog.Info("frans was started in development mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	tracingCleanup, err := otel.NewTracerProvider(ctx, configValue.Otel)
	if err != nil {
		slog.Error("Setup failed", "err", err)
		os.Exit(1)
	}
	defer tracingCleanup()
	r := gin.New()
	r.Use(otelgin.Middleware(otel.TracingService))
	err = routes.SetupRootRouter(r, configValue, db)
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
		startGin(cmd.Context(), configValue, db)
	},
}
