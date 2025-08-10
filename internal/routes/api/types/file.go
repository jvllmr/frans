package apiTypes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/ent"
)

type RequestedFileParam struct {
	ID string `uri:"fileId" binding:"required,uuid"`
}

type PublicFile struct {
	Id              uuid.UUID `json:"id"`
	Sha512          string    `json:"sha512"`
	Size            uint64    `json:"size"`
	Name            string    `json:"name"`
	TimesDownloaded uint64    `json:"timesDownloaded"`
	LastDownloaded  *string   `json:"lastDownloaded"`
}

func ToPublicFile(file *ent.File) PublicFile {
	var lastDownloadedValue *string = nil
	if file.LastDownload != nil {
		formattedValue := file.LastDownload.UTC().Format(http.TimeFormat)
		lastDownloadedValue = &formattedValue
	}

	return PublicFile{
		Id:              file.ID,
		Sha512:          file.Sha512,
		Size:            file.Size,
		Name:            file.Name,
		TimesDownloaded: file.TimesDownloaded,
		LastDownloaded:  lastDownloadedValue,
	}

}
