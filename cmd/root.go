package cmd

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/logging"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "frans",
	Short: "A simple file-sharing tool ready for cloud native",
	Run: func(cmd *cobra.Command, args []string) {
		configValue, db := getConfigAndDBClient()
		defer db.Close()
		go startCronScheduler(configValue, db)
		startGin(configValue, db)
	},
}

func Main() {
	godotenv.Load()
	logConfig, err := config.NewLogConfig()
	if err != nil {
		log.Fatalf("Could not parse logging config: %v", err)
	}
	logging.SetupLogging(logConfig)

	taskCommand.AddCommand(
		sessionLifecycleTaskCommand,
		ticketLifecycleTaskCommand,
		fileLifecycleTaskCommand,
		grantLifecycleTaskCommand,
	)
	rootCmd.AddCommand(taskCommand, cronCmd, serveCmd, migrateCmd)

	cobra.CheckErr(rootCmd.Execute())
}
