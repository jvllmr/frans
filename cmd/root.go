package cmd

import (
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
	"github.com/jvllmr/frans/internal/config"
	"github.com/spf13/cobra"
)

func initFrans(configValue config.Config) {

	config.InitDB(configValue)
	var logHandler slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if configValue.LogJSON {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	}

	basicLogger := slog.New(logHandler)
	slog.SetDefault(basicLogger)

}

var rootCmd = &cobra.Command{
	Use:   "frans",
	Short: "A simple file-sharing tool ready for cloud native",
	Run: func(cmd *cobra.Command, args []string) {
		go startCronScheduler()
		startGin()
	},
}

func Main() {
	godotenv.Load()
	configValue := config.GetSafeConfig()
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
