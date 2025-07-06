package apiRoutes

import (
	"github.com/google/uuid"
	"github.com/jvllmr/frans/pkg/ent"
)

type PublicFile struct {
	Id     uuid.UUID `json:"id"`
	Sha512 string    `json:"sha512"`
	Size   uint64    `json:"size"`
	Name   string    `json:"name"`
}

func ToPublicFile(file *ent.File) PublicFile {
	return PublicFile{
		Id:     file.ID,
		Sha512: file.Sha512,
		Size:   uint64(file.Size),
		Name:   file.Name,
	}

}
