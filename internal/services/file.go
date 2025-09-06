package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
)

type FileService struct {
	config config.Config
}

func (fs FileService) EnsureFilesTmpPath() {
	err := os.MkdirAll(fs.FilesTmpPath(), 0775)
	if err != nil {
		panic(err)
	}
}

func (fs FileService) FilesTmpPath() string {
	return fmt.Sprintf("%s/%s", fs.config.FilesDir, "tmp")
}

func (fs FileService) FilesTmpFilePath() string {
	return fmt.Sprintf("%s/%s", fs.FilesTmpPath(), uuid.New())
}

func (fs FileService) FilesFilePath(fileName string) string {
	return fmt.Sprintf("%s/%s", fs.config.FilesDir, fileName)
}

func (fs FileService) ShouldDeleteFile(

	fileValue *ent.File,
) bool {
	if fileValue.ExpiryType == config.TicketExpiryTypeNone {
		return false
	}
	if fileValue.ExpiryType == config.TicketExpiryTypeSingle {
		return fileValue.TimesDownloaded > 0
	}
	estimatedExpiry := *fs.FileEstimatedExpiry(fileValue)
	now := time.Now()

	if fileValue.ExpiryType == config.TicketExpiryTypeCustom {
		return fileValue.TimesDownloaded >= uint64(fileValue.ExpiryTotalDownloads) ||
			estimatedExpiry.Before(now)
	}
	return fileValue.TimesDownloaded >= uint64(fs.config.DefaultExpiryTotalDownloads) ||
		estimatedExpiry.Before(now)
}

func (fs FileService) DeleteFile(fileValue *ent.File) error {
	filePath := fs.FilesFilePath(fileValue.Sha512)
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	err = config.DBClient.File.DeleteOne(fileValue).Exec(context.Background())
	return err
}

func (fs FileService) FileEstimatedExpiry(fileValue *ent.File) *time.Time {
	return estimatedExpiry(
		fileValue.ExpiryType,
		fs.config.DefaultExpiryTotalDays,
		fs.config.DefaultExpiryDaysSinceLastDownload,
		fileValue.ExpiryTotalDays,
		fileValue.ExpiryDaysSinceLastDownload,
		fileValue.CreatedAt,
		fileValue.LastDownload,
	)
}

type PublicFile struct {
	Id              uuid.UUID `json:"id"`
	Sha512          string    `json:"sha512"`
	Size            uint64    `json:"size"`
	Name            string    `json:"name"`
	CreatedAt       string    `json:"createdAt"`
	TimesDownloaded uint64    `json:"timesDownloaded"`
	LastDownloaded  *string   `json:"lastDownloaded"`
	EstimatedExpiry *string   `json:"estimatedExpiry"`
}

func (fs FileService) ToPublicFile(file *ent.File) PublicFile {
	var lastDownloadedValue *string = nil
	if file.LastDownload != nil {
		formattedValue := file.LastDownload.UTC().Format(http.TimeFormat)
		lastDownloadedValue = &formattedValue
	}

	var estimatedExpiryValue *string = nil

	if estimatedExpiryResult := fs.FileEstimatedExpiry(file); estimatedExpiryResult != nil {
		estimatedExpiry := estimatedExpiryResult.Format(http.TimeFormat)
		estimatedExpiryValue = &estimatedExpiry
	}

	return PublicFile{
		Id:              file.ID,
		Sha512:          file.Sha512,
		Size:            file.Size,
		Name:            file.Name,
		CreatedAt:       file.CreatedAt.UTC().Format(http.TimeFormat),
		TimesDownloaded: file.TimesDownloaded,
		LastDownloaded:  lastDownloadedValue,
		EstimatedExpiry: estimatedExpiryValue,
	}

}

func NewFileService(c config.Config) FileService {
	return FileService{config: c}
}
