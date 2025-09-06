package cmd

import (
	"github.com/jvllmr/frans/internal/config"
	fransCron "github.com/jvllmr/frans/internal/cron"
	"github.com/jvllmr/frans/internal/services"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func startCronScheduler() {
	configValue := config.GetSafeConfig()
	cronRunner := cron.New()
	cronRunner.AddFunc("@every 1m", fransCron.SessionLifecycleTask)

	fs := services.NewFileService(configValue)
	cronRunner.AddFunc("@every 1m", func() {
		fransCron.FileLifecycleTask(fs)
	})

	ts := services.NewTicketService(configValue)
	cronRunner.AddFunc("@every 1m", func() {
		fransCron.TicketsLifecycleTask(ts)
	})

	gs := services.NewGrantService(configValue)
	cronRunner.AddFunc("@every 1m", func() {
		fransCron.GrantsLifecycleTask(gs)
	})

	cronRunner.Run()

}

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start only the cron scheduler",
	Run:   func(cmd *cobra.Command, args []string) { startCronScheduler() },
}
