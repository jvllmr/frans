package routes

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	logging "github.com/jvllmr/frans/internal/logging"
	"github.com/jvllmr/frans/internal/oidc"
	apiRoutes "github.com/jvllmr/frans/internal/routes/api"
	clientRoutes "github.com/jvllmr/frans/internal/routes/client"
)

func SetupRootRouter(r *gin.Engine, configValue config.Config, db *ent.Client) error {

	if err := r.SetTrustedProxies(configValue.TrustedProxies); err != nil {
		return err
	}

	r.Use(logging.GinLogger(slog.Default()), logging.RecoveryLogger(slog.Default()))

	defaultGroup := r.Group(configValue.RootPath)

	oidcProvider, err := oidc.NewOIDC(configValue, db)

	if err != nil {
		return fmt.Errorf("root setup: %w", err)
	}

	clientRoutes.SetupClientRoutes(r, defaultGroup, configValue, db, oidcProvider)
	apiRoutes.SetupAPIRoutes(defaultGroup, configValue, db, oidcProvider)
	return err
}
