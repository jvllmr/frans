package cmd

import (
	"github.com/jvllmr/frans/internal/db"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate frans installation",
	Run:   func(cmd *cobra.Command, args []string) { db.Migrate() },
}
