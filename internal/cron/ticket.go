package fransCron

import (
	"context"
	"log/slog"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/util"
)

func TicketsLifecycleTask(configValue config.Config) {
	tickets := config.DBClient.Ticket.Query().WithFiles().AllX(context.Background())
	deletedCount := 0

	for _, ticketValue := range tickets {
		if util.ShouldDeleteTicket(configValue, ticketValue) {
			err := config.DBClient.Ticket.DeleteOne(ticketValue).Exec(context.Background())
			if err != nil {
				slog.Error("Could not delete ticket", "err", err)
				continue
			}
			deletedCount += 1
		}

	}
	slog.Info("Deleted tickets", "count", deletedCount)
}
