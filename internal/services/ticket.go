package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
)

type TicketService struct {
	config config.Config
}

func (ts TicketService) TicketEstimatedExpiry(ticketValue *ent.Ticket) *time.Time {
	var latestDownload *time.Time = nil
	for _, file := range ticketValue.Edges.Files {
		if latestDownload == nil ||
			(file.LastDownload != nil && latestDownload.Before(*file.LastDownload)) {
			latestDownload = file.LastDownload
		}
	}

	return estimatedExpiry(
		ticketValue.ExpiryType,
		ts.config.DefaultExpiryTotalDays,
		ts.config.DefaultExpiryDaysSinceLastDownload,
		ticketValue.ExpiryTotalDays,
		ticketValue.ExpiryDaysSinceLastDownload,
		ticketValue.CreatedAt,
		latestDownload,
	)
}

func (ts TicketService) TicketShareLink(ctx *gin.Context, ticket *ent.Ticket) string {
	return fmt.Sprintf("%s/s/%s", ts.config.GetBaseURL(ctx.Request), ticket.ID.String())
}

func (ts TicketService) ShouldDeleteTicket(ticketValue *ent.Ticket) bool {
	estimatedExpiry := ts.TicketEstimatedExpiry(ticketValue)
	now := time.Now()

	return len(ticketValue.Edges.Files) == 0 ||
		(estimatedExpiry != nil && estimatedExpiry.Before(now))
}

type PublicTicket struct {
	Id              uuid.UUID    `json:"id"`
	Comment         *string      `json:"comment"`
	EstimatedExpiry *string      `json:"estimatedExpiry"`
	User            PublicUser   `json:"owner"`
	Files           []PublicFile `json:"files"`
	CreatedAt       string       `json:"createdAt"`
}

func (ts TicketService) ToPublicTicket(fs FileService, ticket *ent.Ticket) PublicTicket {
	files := make([]PublicFile, len(ticket.Edges.Files))
	for i, file := range ticket.Edges.Files {
		files[i] = fs.ToPublicFile(file)
	}
	var estimatedExpiryValue *string = nil

	if estimatedExpiryResult := ts.TicketEstimatedExpiry(ticket); estimatedExpiryResult != nil {
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

func NewTicketService(c config.Config) TicketService {
	return TicketService{config: c}
}
