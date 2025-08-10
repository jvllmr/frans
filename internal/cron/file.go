package fransCron

import (
	"context"
	"log/slog"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/file"
	"github.com/jvllmr/frans/internal/util"
)

func FileLifecycleTask(configValue config.Config) {
	files := config.DBClient.File.Query().
		Where(file.TimesDownloadedGT(0)).
		WithTickets(func(ticketQuery *ent.TicketQuery) {
			ticketQuery.WithOwner()
		}).
		AllX(context.Background())
	deletedCount := 0
	var users []*ent.User
	for _, fileValue := range files {
		var ticketValue *ent.Ticket
		if len(fileValue.Edges.Tickets) == 1 {
			ticketValue = fileValue.Edges.Tickets[0]
		}

		if ticketValue == nil || util.ShouldDeleteFileConnectedToTicket(configValue,
			*ticketValue, *fileValue) {
			err := util.DeleteFile(configValue, fileValue)
			if err != nil {
				filePath := util.GetFilesFilePath(configValue, fileValue.Sha512)
				slog.Error("Could not delete file", "file", filePath, "err", err)
				continue
			}

			deletedCount += 1
			if ticketValue != nil {
				users = append(users, ticketValue.Edges.Owner)
			}
		}

	}
	slog.Info("Deleted files", "count", deletedCount)

	for _, userValue := range users {
		util.RefreshUserTotalDataSize(context.Background(), userValue)
	}
	slog.Info(
		"Refreshed totalDataSize field for all users affected by file deletions",
		"count",
		len(users),
	)
}
