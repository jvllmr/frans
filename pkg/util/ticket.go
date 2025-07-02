package util

import (
	"time"

	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent"
)

func GetEstimatedExpiry(configValue config.Config, ticket *ent.Ticket) time.Time {
	expiryTotalDays := configValue.DefaultExpiryTotalDays
	expiryDaysSinceLastDownload := configValue.DefaultExpiryDaysSinceLastDownload
	if ticket.ExpiryType == "custom" {
		expiryTotalDays = ticket.ExpiryTotalDays
		expiryDaysSinceLastDownload = ticket.ExpiryDaysSinceLastDownload
	}

	totalLimit := ticket.CreatedAt.Add(time.Hour * 24 * time.Duration(expiryTotalDays))
	lastDownloadLimit := ticket.LastDownload.Add(
		time.Hour * 24 * time.Duration(expiryDaysSinceLastDownload),
	)

	if lastDownloadLimit.Before(totalLimit) {
		return lastDownloadLimit.UTC()
	}

	return totalLimit.UTC()
}
