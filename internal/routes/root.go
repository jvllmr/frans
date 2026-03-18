package routes

import (
	"fmt"
	"log/slog"

	"codeberg.org/jvllmr/frans/internal/config"
	"codeberg.org/jvllmr/frans/internal/ent"
	logging "codeberg.org/jvllmr/frans/internal/logging"
	"codeberg.org/jvllmr/frans/internal/oidc"
	apiRoutes "codeberg.org/jvllmr/frans/internal/routes/api"
	clientRoutes "codeberg.org/jvllmr/frans/internal/routes/client"
	"github.com/gin-gonic/gin"
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
