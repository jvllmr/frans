package main

import (
	"fmt"
	"log/slog"
	"os"

	ginlogger "github.com/FabienMht/ginslog/logger"
	ginrecovery "github.com/FabienMht/ginslog/recovery"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/jvllmr/frans/pkg/config"
	apiRoutes "github.com/jvllmr/frans/pkg/routes/api"
	clientRoutes "github.com/jvllmr/frans/pkg/routes/client"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

func setupRouter(configValue config.Config) *gin.Engine {

	r := gin.New()

	var stdoutHandler slog.Handler
	if configValue.LogJSON {
		stdoutHandler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		stdoutHandler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(slogmulti.Fanout(otelslog.NewHandler("frans"), stdoutHandler))
	slog.SetDefault(logger)
	r.Use(ginlogger.New(logger), ginrecovery.New(logger))

	defaultGroup := r.Group(configValue.RootPath)
	clientRoutes.SetupClientRoutes(r, defaultGroup, configValue)
	apiRoutes.SetupAPIRoutes(defaultGroup, configValue)
	return r
}

func main() {
	godotenv.Load()
	configValue, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	config.InitDB(configValue)
	defer config.DBClient.Close()
	config.InitOIDC(configValue)
	if configValue.DevMode {
		slog.Info("frans was started in development mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := setupRouter(configValue)
	serveString := fmt.Sprintf("%s:%d", configValue.Host, configValue.Port)
	slog.Info(fmt.Sprintf("Serving on %s", serveString))
	r.Run(serveString)

}
