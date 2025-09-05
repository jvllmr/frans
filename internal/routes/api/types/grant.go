package apiTypes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/util"
)

type PublicGrant struct {
	Id              uuid.UUID    `json:"id"`
	Comment         *string      `json:"comment"`
	EstimatedExpiry *string      `json:"estimatedExpiry"`
	User            PublicUser   `json:"owner"`
	Files           []PublicFile `json:"files"`
	CreatedAt       string       `json:"createdAt"`
}

func ToPublicGrant(
	configValue config.Config,
	grantValue *ent.Grant,
	files []*ent.File,
) PublicGrant {
	publicFiles := make([]PublicFile, len(files))
	for i, file := range files {
		publicFiles[i] = ToPublicFile(configValue, file)
	}
	var estimatedExpiryValue *string = nil

	if estimatedExpiryResult := util.GrantEstimatedExpiry(configValue, grantValue); estimatedExpiryResult != nil {
		estimatedExpiry := estimatedExpiryResult.Format(http.TimeFormat)
		estimatedExpiryValue = &estimatedExpiry
	}

	return PublicGrant{
		Id:              grantValue.ID,
		Comment:         grantValue.Comment,
		User:            ToPublicUser(grantValue.Edges.Owner),
		EstimatedExpiry: estimatedExpiryValue,
		CreatedAt:       grantValue.CreatedAt.UTC().Format(http.TimeFormat),
		Files:           publicFiles,
	}
}
