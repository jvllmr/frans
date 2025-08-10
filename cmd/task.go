package cmd

import "github.com/spf13/cobra"

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
