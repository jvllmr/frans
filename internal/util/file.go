package util

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/grant"
	"github.com/jvllmr/frans/internal/ent/ticket"
	"github.com/jvllmr/frans/internal/ent/user"
)

func EnsureFilesTmpPath(configValue config.Config) {
	err := os.MkdirAll(FilesTmpPath(configValue), 0775)
	if err != nil {
		panic(err)
	}
}

func FilesTmpPath(configValue config.Config) string {
	return fmt.Sprintf("%s/%s", configValue.FilesDir, "tmp")
}

func FilesTmpFilePath(configValue config.Config) string {
	return fmt.Sprintf("%s/%s", FilesTmpPath(configValue), uuid.New())
}

func FilesFilePath(configValue config.Config, fileName string) string {
	return fmt.Sprintf("%s/%s", configValue.FilesDir, fileName)
}

func ShouldDeleteFile(
	configValue config.Config,
	fileValue *ent.File,
) bool {
	if fileValue.ExpiryType == config.TicketExpiryTypeNone {
		return false
	}
	if fileValue.ExpiryType == config.TicketExpiryTypeSingle {
		return fileValue.TimesDownloaded > 0
	}
	estimatedExpiry := *FileEstimatedExpiry(configValue, fileValue)
	now := time.Now()

	if fileValue.ExpiryType == config.TicketExpiryTypeCustom {
		return fileValue.TimesDownloaded >= uint64(fileValue.ExpiryTotalDownloads) ||
			estimatedExpiry.Before(now)
	}
	return fileValue.TimesDownloaded >= uint64(configValue.DefaultExpiryTotalDownloads) ||
		estimatedExpiry.Before(now)
}

func DeleteFile(configValue config.Config, fileValue *ent.File) error {
	filePath := FilesFilePath(configValue, fileValue.Sha512)
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	err = config.DBClient.File.DeleteOne(fileValue).Exec(context.Background())
	return err
}

func UserHasFileAccess(ctx context.Context, userValue *ent.User, fileValue *ent.File) bool {
	return userValue.IsAdmin || fileValue.QueryTickets().
		Where(ticket.HasOwnerWith(user.ID(userValue.ID))).
		CountX(ctx) > 0 || fileValue.QueryGrants().
		Where(grant.HasOwnerWith(user.ID(userValue.ID))).
		CountX(ctx) > 0
}

func FileEstimatedExpiry(configValue config.Config, fileValue *ent.File) *time.Time {
	return estimatedExpiry(
		fileValue.ExpiryType,
		configValue.DefaultExpiryTotalDays,
		configValue.DefaultExpiryDaysSinceLastDownload,
		fileValue.ExpiryTotalDays,
		fileValue.ExpiryDaysSinceLastDownload,
		fileValue.CreatedAt,
		fileValue.LastDownload,
	)
}
