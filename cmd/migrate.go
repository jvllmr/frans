package cmd

import (
	"github.com/jvllmr/frans/internal/migration"
	"github.com/spf13/cobra"
)

func migrateFrans() {
	migration.Migrate()
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate frans installation",
	Run:   func(cmd *cobra.Command, args []string) { migrateFrans() },
}
