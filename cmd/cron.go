package cmd

import (
	"log"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/services"
	fransCron "github.com/jvllmr/frans/internal/tasks"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func startCronScheduler(configValue config.Config, db *ent.Client) {

	cronRunner := cron.New()

	_, err := cronRunner.AddFunc("@every 1m", func() { fransCron.SessionLifecycleTask(db) })
	if err != nil {
		log.Fatalf("create sessions cronjob: %v", err)
	}

	fs := services.NewFileService(configValue, db)
	_, err = cronRunner.AddFunc("@every 1m", func() {
		fransCron.FileLifecycleTask(db, fs)
	})

	if err != nil {
		log.Fatalf("create files cronjob: %v", err)
	}

	ts := services.NewTicketService(configValue, db)
	_, err = cronRunner.AddFunc("@every 1m", func() {
		fransCron.TicketsLifecycleTask(db, ts)
	})

	if err != nil {
		log.Fatalf("create tickets cronjob: %v", err)
	}

	gs := services.NewGrantService(configValue)
	_, err = cronRunner.AddFunc("@every 1m", func() {
		fransCron.GrantsLifecycleTask(db, gs)
	})

	if err != nil {
		log.Fatalf("create grants cronjob: %v", err)
	}

	cronRunner.Run()

}

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start only the cron scheduler",
	Run: func(cmd *cobra.Command, args []string) {
		configValue, db := getConfigAndDBClient()
		defer func() {
			if err := db.Close(); err != nil {
				log.Fatalf("could not close db connection: %v", err)
			}
		}()
		startCronScheduler(configValue, db)
	},
}
