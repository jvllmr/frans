package routes

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	logging "github.com/jvllmr/frans/internal/logging"
	apiRoutes "github.com/jvllmr/frans/internal/routes/api"
	clientRoutes "github.com/jvllmr/frans/internal/routes/client"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

func SetupRootRouter(configValue config.Config) *gin.Engine {

	r := gin.New()

	var stdoutHandler slog.Handler
	if configValue.LogJSON {
		stdoutHandler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		stdoutHandler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(slogmulti.Fanout(otelslog.NewHandler("frans"), stdoutHandler))
	slog.SetDefault(logger)
	r.Use(logging.GinLogger(logger), logging.RecoveryLogger(logger))

	defaultGroup := r.Group(configValue.RootPath)
	clientRoutes.SetupClientRoutes(r, defaultGroup, configValue)
	apiRoutes.SetupAPIRoutes(defaultGroup, configValue)
	return r
}
