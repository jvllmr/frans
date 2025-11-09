package cmd

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jvllmr/frans/internal/db/sqlite3"

	"github.com/joho/godotenv"
	"github.com/jvllmr/frans/internal/db"
	"github.com/jvllmr/frans/internal/logging"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "frans",
	Short: "A simple file-sharing tool ready for cloud native",
	Run: func(cmd *cobra.Command, args []string) {
		configValue, dbCon := getConfigAndDBClient()
		defer func() {
			if err := dbCon.Close(); err != nil {
				log.Fatalf("could not close db connection: %v", err)
			}
		}()

		if !configValue.DevMode {
			db.Migrate(configValue.DBConfig)
		}
		go startCronScheduler(configValue, dbCon)
		startGin(cmd.Context(), configValue, dbCon)
	},
}

func Main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("loading vars from .env: %v", err)
	}
	err = logging.SetupLogging()
	if err != nil {
		log.Fatalf("could not setup logging: %v", err)
	}

	taskCommand.AddCommand(
		sessionLifecycleTaskCommand,
		ticketLifecycleTaskCommand,
		fileLifecycleTaskCommand,
		grantLifecycleTaskCommand,
	)
	rootCmd.AddCommand(taskCommand, cronCmd, serveCmd, migrateCmd)

	cobra.CheckErr(rootCmd.Execute())
}
