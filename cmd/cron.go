package cmd

import (
	"github.com/jvllmr/frans/internal/config"
	fransCron "github.com/jvllmr/frans/internal/cron"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/services"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func startCronScheduler(configValue config.Config, db *ent.Client) {

	cronRunner := cron.New()

	cronRunner.AddFunc("@every 1m", func() { fransCron.SessionLifecycleTask(db) })

	fs := services.NewFileService(configValue, db)
	cronRunner.AddFunc("@every 1m", func() {
		fransCron.FileLifecycleTask(db, fs)
	})

	ts := services.NewTicketService(configValue)
	cronRunner.AddFunc("@every 1m", func() {
		fransCron.TicketsLifecycleTask(db, ts)
	})

	gs := services.NewGrantService(configValue)
	cronRunner.AddFunc("@every 1m", func() {
		fransCron.GrantsLifecycleTask(db, gs)
	})

	cronRunner.Run()

}

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start only the cron scheduler",
	Run: func(cmd *cobra.Command, args []string) {
		configValue, db := getConfigAndDBClient()
		startCronScheduler(configValue, db)
	},
}
