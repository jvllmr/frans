package routes

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	logging "github.com/jvllmr/frans/internal/logging"
	"github.com/jvllmr/frans/internal/oidc"
	apiRoutes "github.com/jvllmr/frans/internal/routes/api"
	clientRoutes "github.com/jvllmr/frans/internal/routes/client"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

func SetupRootRouter(configValue config.Config, db *ent.Client) (*gin.Engine, error) {

	r := gin.New()

	r.SetTrustedProxies(configValue.TrustedProxies)
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

	oidcProvider, err := oidc.NewOIDC(configValue, db)

	if err != nil {
		return nil, fmt.Errorf("root setup: %w", err)
	}

	clientRoutes.SetupClientRoutes(r, defaultGroup, configValue, db, oidcProvider)
	apiRoutes.SetupAPIRoutes(defaultGroup, configValue, db, oidcProvider)
	return r, nil
}
