package util

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
)

func GrantShareLink(ctx *gin.Context, configValue config.Config, grant *ent.Grant) string {
	return fmt.Sprintf("%s/s/%s", config.GetBaseURL(configValue, ctx.Request), grant.ID.String())
}

func GrantEstimatedExpiry(configValue config.Config, grantValue *ent.Grant) *time.Time {
	return estimatedExpiry(
		grantValue.ExpiryType,
		configValue.GrantDefaultExpiryTotalDays,
		configValue.GrantDefaultExpiryDaysSinceLastUpload,
		grantValue.ExpiryTotalDays,
		grantValue.ExpiryDaysSinceLastUpload,
		grantValue.CreatedAt,
		grantValue.LastUpload,
	)
}

func GrantExpired(configValue config.Config, grantValue *ent.Grant) bool {
	if grantValue.ExpiryType == config.TicketExpiryTypeNone {
		return false
	}
	if grantValue.ExpiryType == config.TicketExpiryTypeSingle {
		return grantValue.TimesUploaded > 0
	}
	estimatedExpiryValue := *GrantEstimatedExpiry(configValue, grantValue)
	now := time.Now()

	if grantValue.ExpiryType == config.TicketExpiryTypeCustom {
		return grantValue.TimesUploaded >= uint64(grantValue.ExpiryTotalUploads) ||
			estimatedExpiryValue.Before(now)
	}

	return grantValue.TimesUploaded >= uint64(configValue.GrantDefaultExpiryTotalUploads) ||
		estimatedExpiryValue.Before(now)

}

func ShouldDeleteGrant(configValue config.Config, grantValue *ent.Grant) bool {
	return len(grantValue.Edges.Files) == 0 && GrantExpired(configValue, grantValue)
}
