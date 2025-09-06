package services

import "time"

func estimatedExpiry(
	expiryType string,
	defaultExpiryTotalDays uint8,
	defaultExpiryDaysSinceLastDownload uint8,
	customExpiryTotalDays uint8,
	customExpiryDaysSinceLastDownload uint8,
	createdAt time.Time,
	lastDownload *time.Time,
) *time.Time {
	if expiryType == "none" {
		return nil
	}

	expiryTotalDays := defaultExpiryTotalDays
	expiryDaysSinceLastDownload := defaultExpiryDaysSinceLastDownload
	if expiryType == "custom" {
		expiryTotalDays = customExpiryTotalDays
		expiryDaysSinceLastDownload = customExpiryDaysSinceLastDownload
	}

	totalLimit := createdAt.Add(time.Hour * 24 * time.Duration(expiryTotalDays)).UTC()

	if lastDownload == nil {
		return &totalLimit
	}

	lastDownloadLimit := lastDownload.Add(
		time.Hour * 24 * time.Duration(expiryDaysSinceLastDownload),
	).UTC()

	if lastDownloadLimit.Before(totalLimit) {
		return &lastDownloadLimit
	}

	return &totalLimit
}
