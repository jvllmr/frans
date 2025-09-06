package tasks

import (
	"context"
	"log/slog"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/file"
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/util"
)

func FileLifecycleTask(db *ent.Client, fs services.FileService) {
	files := db.File.Query().
		Where(file.TimesDownloadedGT(0)).
		WithTickets(func(ticketQuery *ent.TicketQuery) {
			ticketQuery.WithOwner()
		}).
		WithGrants(func(grantQuery *ent.GrantQuery) {
			grantQuery.WithOwner()
		}).
		AllX(context.Background())
	deletedCount := 0
	var users []*ent.User
	for _, fileValue := range files {
		var ticketValue *ent.Ticket
		if len(fileValue.Edges.Tickets) == 1 {
			ticketValue = fileValue.Edges.Tickets[0]
		}

		if fs.ShouldDeleteFile(fileValue) {
			err := fs.DeleteFile(fileValue)
			if err != nil {
				filePath := fs.FilesFilePath(fileValue.Sha512)
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
		"Refreshed totalDataSize field for all users affected by file deletions",
		"count",
		len(users),
	)
}
