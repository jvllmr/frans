package cmd

import (
	"github.com/jvllmr/frans/internal/services"
	fransCron "github.com/jvllmr/frans/internal/tasks"
	"github.com/spf13/cobra"
)

var taskCommand = &cobra.Command{
	Use:   "task",
	Short: "Start specific cron task once",
}

var sessionLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-session",
	Short: "Delete expired sessions",
	Run: func(cmd *cobra.Command, args []string) {
		_, db := getConfigAndDBClient()
		defer db.Close()
		fransCron.SessionLifecycleTask(db)
	},
}

var ticketLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-ticket",
	Short: "Delete expired tickets",
	Run: func(cmd *cobra.Command, args []string) {
		configValue, db := getConfigAndDBClient()
		defer db.Close()
		ts := services.NewTicketService(configValue)
		fransCron.TicketsLifecycleTask(db, ts)
	},
}

var grantLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-grant",
	Short: "Delete expired grants",
	Run: func(cmd *cobra.Command, args []string) {
		configValue, db := getConfigAndDBClient()
		defer db.Close()
		gs := services.NewGrantService(configValue)
		fransCron.GrantsLifecycleTask(db, gs)
	},
}

var fileLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-file",
	Short: "Delete expired files",
	Run: func(cmd *cobra.Command, args []string) {
		configValue, db := getConfigAndDBClient()
		defer db.Close()
		fs := services.NewFileService(configValue, db)
		fransCron.FileLifecycleTask(db, fs)
	},
}
