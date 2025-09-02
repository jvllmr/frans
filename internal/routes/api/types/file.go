package apiTypes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/util"
)

type RequestedFileParam struct {
	ID string `uri:"fileId" binding:"required,uuid"`
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

func ToPublicFile(configValue config.Config, file *ent.File) PublicFile {
	var lastDownloadedValue *string = nil
	if file.LastDownload != nil {
		formattedValue := file.LastDownload.UTC().Format(http.TimeFormat)
		lastDownloadedValue = &formattedValue
	}

	var estimatedExpiryValue *string = nil

	if estimatedExpiryResult := util.FileEstimatedExpiry(configValue, file); estimatedExpiryResult != nil {
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
