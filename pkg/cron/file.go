package fransCron

import (
	"context"
	"log/slog"
	"os"

	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent/file"
	"github.com/jvllmr/frans/pkg/util"
)

func FileLifecycleTask(configValue config.Config) {
	files := config.DBClient.File.Query().
		Where(file.TimesDownloadedGT(0)).
		WithTickets().
		AllX(context.Background())
	deletedCount := 0
	for _, fileValue := range files {
		skip := false
		for _, ticketValue := range fileValue.Edges.Tickets {
			if ticketValue.ExpiryType == config.TicketExpiryTypeNone {
				continue
			}
			if ticketValue.ExpiryType == config.TicketExpiryTypeCustom &&
				ticketValue.ExpiryTotalDownloads > uint8(fileValue.TimesDownloaded) ||
				ticketValue.ExpiryType == config.TicketExpiryTypeAuto &&
					configValue.DefaultExpiryTotalDownloads > uint8(fileValue.TimesDownloaded) {
				skip = true
			}
		}
		if !skip {
			filePath := util.GetFilesFilePath(configValue, fileValue.Sha512)
			err := os.Remove(filePath)
			if err != nil {
				slog.Error("Could not delete file", "file", filePath, "err", err)
				continue
			}
			config.DBClient.File.DeleteOne(fileValue).ExecX(context.Background())
			deletedCount += 1
		}
	}
	slog.Info("Deleted files", "count", deletedCount)
}
