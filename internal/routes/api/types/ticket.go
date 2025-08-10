package apiTypes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/util"
)

type PublicTicket struct {
	Id              uuid.UUID    `json:"id"`
	Comment         *string      `json:"comment"`
	EstimatedExpiry *string      `json:"estimatedExpiry"`
	User            PublicUser   `json:"owner"`
	Files           []PublicFile `json:"files"`
	CreatedAt       string       `json:"createdAt"`
}

func ToPublicTicket(configValue config.Config, ticket *ent.Ticket) PublicTicket {
	files := []PublicFile{}
	for _, file := range ticket.Edges.Files {
		files = append(files, ToPublicFile(file))
	}
	var estimatedExpiryValue *string = nil

	if estimatedExpiryResult := util.GetEstimatedExpiry(configValue, ticket); estimatedExpiryResult != nil {
		estimatedExpiry := estimatedExpiryResult.Format(http.TimeFormat)
		estimatedExpiryValue = &estimatedExpiry
	}

	return PublicTicket{
		Id:              ticket.ID,
		Comment:         ticket.Comment,
		User:            ToPublicUser(ticket.Edges.Owner),
		EstimatedExpiry: estimatedExpiryValue,
		Files:           files,
		CreatedAt:       ticket.CreatedAt.UTC().Format(http.TimeFormat),
	}
}
