package util

import (
	"time"

	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent"
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
