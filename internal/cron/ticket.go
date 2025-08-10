package fransCron

import (
	"context"
	"log/slog"
	"time"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/util"
)

func TicketsLifecycleTask(configValue config.Config) {
	tickets := config.DBClient.Ticket.Query().WithFiles().AllX(context.Background())
	deletedCount := 0
	for _, ticketValue := range tickets {
		estimatedExpiry := util.GetEstimatedExpiry(configValue, ticketValue)
		now := time.Now()
		if len(ticketValue.Edges.Files) > 0 &&
			(estimatedExpiry == nil || estimatedExpiry.After(now)) {
			continue
		}
		config.DBClient.Ticket.DeleteOne(ticketValue).ExecX(context.Background())
		deletedCount += 1
	}
	slog.Info("Deleted tickets", "count", deletedCount)
}
