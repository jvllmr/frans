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
		WithData().
		WithOwner().
		AllX(context.Background())
	deletedCount := 0
	var users []*ent.User
	for _, fileValue := range files {
		if fs.ShouldDeleteFile(fileValue) {
			fileOwner := fileValue.Edges.Owner
			err := fs.DeleteFile(context.Background(), fileValue)
			if err != nil {
				filePath := fs.FilesFilePath(fileValue.Edges.Data.ID)
				slog.Error("Could not delete file", "file", filePath, "err", err)
				continue
			}
			deletedCount += 1
			users = append(users, fileOwner)

		}

	}
	slog.Info("Deleted files", "count", deletedCount)
	tx, err := db.Tx(context.Background())
	if err != nil {
		slog.Error("Could not update users data size", "err", err)
		return
	}
	for _, u := range users {
		err := util.RefreshUserTotalDataSize(context.Background(), u, tx)
		if err != nil {
			slog.Error(
				"Could not refresh total data size of user",
				"err",
				err,
				"user",
				u.Username,
			)
		}
	}
	if err := tx.Commit(); err != nil {
		slog.Error("Could not commit users data size updates", "err", err)
		return
	}
	slog.Info(
		"Refreshed totalDataSize field for all users affected by file deletions",
		"count",
		len(users),
	)
}
