package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/jvllmr/frans/pkg/config"
	fransCron "github.com/jvllmr/frans/pkg/cron"
	logging "github.com/jvllmr/frans/pkg/logger"
	apiRoutes "github.com/jvllmr/frans/pkg/routes/api"
	clientRoutes "github.com/jvllmr/frans/pkg/routes/client"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/cobra"

	"github.com/robfig/cron/v3"
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
	r.Use(logging.GinLogger(logger), logging.RecoveryLogger(logger))

	defaultGroup := r.Group(configValue.RootPath)
	clientRoutes.SetupClientRoutes(r, defaultGroup, configValue)
	apiRoutes.SetupAPIRoutes(defaultGroup, configValue)
	return r
}

func getSafeConfig() config.Config {
	configValue, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	return configValue
}

func initFrans(configValue config.Config) {

	config.InitDB(configValue)
	var logHandler slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if configValue.LogJSON {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	}

	basicLogger := slog.New(logHandler)
	slog.SetDefault(basicLogger)

}

func startGin() {
	configValue := getSafeConfig()
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

func startCronScheduler() {
	configValue := getSafeConfig()
	cronRunner := cron.New()
	cronRunner.AddFunc("@every 1m", fransCron.SessionLifecycleTask)
	cronRunner.AddFunc("@every 1m", func() {
		fransCron.FileLifecycleTask(configValue)
	})

	cronRunner.AddFunc("@every 1m", func() {
		fransCron.TicketsLifecycleTask(configValue)
	})

	cronRunner.Run()

}

var rootCmd = &cobra.Command{
	Use:   "frans",
	Short: "A simple file-sharing tool ready for cloud native",
	Run: func(cmd *cobra.Command, args []string) {
		go startCronScheduler()
		startGin()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start only the server",
	Run:   func(cmd *cobra.Command, args []string) { startGin() },
}

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start only the cron scheduler",
	Run:   func(cmd *cobra.Command, args []string) { startCronScheduler() },
}

var taskCommand = &cobra.Command{
	Use:   "task",
	Short: "Start specific cron task once",
}

var sessionLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-session",
	Short: "Delete stale sessions",
}

var ticketLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-ticket",
	Short: "Delete stale tickets",
}

var fileLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-file",
	Short: "Delete stale files",
}

func main() {
	godotenv.Load()
	configValue := getSafeConfig()
	initFrans(configValue)
	defer config.DBClient.Close()

	taskCommand.AddCommand(
		sessionLifecycleTaskCommand,
		ticketLifecycleTaskCommand,
		fileLifecycleTaskCommand,
	)
	rootCmd.AddCommand(taskCommand, cronCmd, serveCmd)

	cobra.CheckErr(rootCmd.Execute())
}
