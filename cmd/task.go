package cmd

import "github.com/spf13/cobra"

var taskCommand = &cobra.Command{
	Use:   "task",
	Short: "Start specific cron task once",
}

var sessionLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-session",
	Short: "Delete expired sessions",
}

var ticketLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-ticket",
	Short: "Delete expired tickets",
}

var grantLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-grant",
	Short: "Delete expired grants",
}

var fileLifecycleTaskCommand = &cobra.Command{
	Use:   "lifecycle-file",
	Short: "Delete expired files",
}
