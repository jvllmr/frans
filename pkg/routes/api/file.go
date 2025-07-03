package apiRoutes

import "github.com/jvllmr/frans/pkg/ent"

type PublicFile struct {
	Sha256 string `json:"sha256sum"`
	Size   uint64 `json:"size"`
	Name   string `json:"name"`
}

func ToPublicFile(file *ent.File) PublicFile {
	return PublicFile{
		Sha256: file.ID,
		Size:   uint64(file.Size),
		Name:   file.Name,
	}

}
