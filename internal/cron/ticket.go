package fransCron

import (
	"context"
	"log/slog"
	"time"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/util"
)

func TicketsLifecycleTask(configValue config.Config) {
	tickets := config.DBClient.Ticket.Query().WithFiles().WithOwner().AllX(context.Background())
	deletedCount := 0
	var users []*ent.User
	for _, ticketValue := range tickets {
		estimatedExpiry := util.GetEstimatedExpiry(configValue, ticketValue)
		now := time.Now()
		if len(ticketValue.Edges.Files) > 0 &&
			(estimatedExpiry == nil || estimatedExpiry.After(now)) {
			continue
		}
		config.DBClient.Ticket.DeleteOne(ticketValue).ExecX(context.Background())
		deletedCount += 1
		users = append(users, ticketValue.Edges.Owner)
	}
	slog.Info("Deleted tickets", "count", deletedCount)

	for _, userValue := range users {
		err := util.RefreshUserTotalDataSize(context.Background(), userValue)
		if err != nil {
			slog.Error(
				"Could not refresh total data size of user",
				"err",
				err,
				"user",
				userValue.Username,
			)
		}
	}
	slog.Info(
		"Refreshed totalDataSize field for all users affected by ticket deletions",
		"count",
		len(users),
	)
}
