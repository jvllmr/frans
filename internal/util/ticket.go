package util

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
)

func GetEstimatedExpiry(configValue config.Config, ticket *ent.Ticket) *time.Time {
	if ticket.ExpiryType == "none" {
		return nil
	}

	expiryTotalDays := configValue.DefaultExpiryTotalDays
	expiryDaysSinceLastDownload := configValue.DefaultExpiryDaysSinceLastDownload
	if ticket.ExpiryType == "custom" {
		expiryTotalDays = ticket.ExpiryTotalDays
		expiryDaysSinceLastDownload = ticket.ExpiryDaysSinceLastDownload
	}

	totalLimit := ticket.CreatedAt.Add(time.Hour * 24 * time.Duration(expiryTotalDays)).UTC()
	var latestDownload *time.Time = nil
	for _, file := range ticket.Edges.Files {
		if latestDownload == nil ||
			(file.LastDownload != nil && latestDownload.Before(*file.LastDownload)) {
			latestDownload = file.LastDownload
		}
	}
	if latestDownload == nil {
		return &totalLimit
	}

	lastDownloadLimit := latestDownload.Add(
		time.Hour * 24 * time.Duration(expiryDaysSinceLastDownload),
	).UTC()

	if lastDownloadLimit.Before(totalLimit) {
		return &lastDownloadLimit
	}

	return &totalLimit
}

func GetTicketShareLink(ctx *gin.Context, configValue config.Config, ticket *ent.Ticket) string {
	return fmt.Sprintf("%s/s/%s", config.GetBaseURL(configValue, ctx.Request), ticket.ID.String())
}
