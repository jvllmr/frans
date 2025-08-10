package util

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
)

func EnsureFilesTmpPath(configValue config.Config) {
	err := os.MkdirAll(GetFilesTmpPath(configValue), 0775)
	if err != nil {
		panic(err)
	}
}

func GetFilesTmpPath(configValue config.Config) string {
	return fmt.Sprintf("%s/%s", configValue.FilesDir, "tmp")
}

func GetFilesTmpFilePath(configValue config.Config) string {
	return fmt.Sprintf("%s/%s", GetFilesTmpPath(configValue), uuid.New())
}

func GetFilesFilePath(configValue config.Config, fileName string) string {
	return fmt.Sprintf("%s/%s", configValue.FilesDir, fileName)
}

func ShouldDeleteFileConnectedToTicket(
	configValue config.Config,
	ticketValue ent.Ticket,
	fileValue ent.File,
) bool {
	return ticketValue.ExpiryType != config.TicketExpiryTypeNone ||
		(ticketValue.ExpiryType == config.TicketExpiryTypeCustom &&
			ticketValue.ExpiryTotalDownloads > uint8(fileValue.TimesDownloaded) ||
			ticketValue.ExpiryType == config.TicketExpiryTypeAuto &&
				configValue.DefaultExpiryTotalDownloads > uint8(fileValue.TimesDownloaded))
}

func DeleteFile(configValue config.Config, fileValue *ent.File) error {
	filePath := GetFilesFilePath(configValue, fileValue.Sha512)
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	err = config.DBClient.File.DeleteOne(fileValue).Exec(context.Background())
	return err
}
