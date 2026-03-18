package cmd

import (
	"log"

	"codeberg.org/jvllmr/frans/internal/config"
	"codeberg.org/jvllmr/frans/internal/db"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate frans installation",
	Run: func(cmd *cobra.Command, args []string) {
		dbConfig, err := config.NewDBConfig()
		if err != nil {
			log.Fatalf("Could not get database config: %v", err)
		}
		db.Migrate(dbConfig.DBConfig)
	},
}
