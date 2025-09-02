package util

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
)

func TicketEstimatedExpiry(configValue config.Config, ticketValue *ent.Ticket) *time.Time {
	var latestDownload *time.Time = nil
	for _, file := range ticketValue.Edges.Files {
		if latestDownload == nil ||
			(file.LastDownload != nil && latestDownload.Before(*file.LastDownload)) {
			latestDownload = file.LastDownload
		}
	}

	return estimatedExpiry(
		ticketValue.ExpiryType,
		configValue.DefaultExpiryTotalDays,
		configValue.DefaultExpiryDaysSinceLastDownload,
		ticketValue.ExpiryTotalDays,
		ticketValue.ExpiryDaysSinceLastDownload,
		ticketValue.CreatedAt,
		latestDownload,
	)
}

func TicketShareLink(ctx *gin.Context, configValue config.Config, ticket *ent.Ticket) string {
	return fmt.Sprintf("%s/s/%s", config.GetBaseURL(configValue, ctx.Request), ticket.ID.String())
}

func ShouldDeleteTicket(configValue config.Config, ticketValue *ent.Ticket) bool {
	estimatedExpiry := TicketEstimatedExpiry(configValue, ticketValue)
	now := time.Now()

	return len(ticketValue.Edges.Files) == 0 ||
		(estimatedExpiry != nil && estimatedExpiry.Before(now))
}
