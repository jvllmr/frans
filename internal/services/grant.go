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

type GrantService struct {
	config config.Config
}

func (gs GrantService) GrantShareLink(ctx *gin.Context, grant *ent.Grant) string {
	return fmt.Sprintf("%s/s/%s", config.GetBaseURL(gs.config, ctx.Request), grant.ID.String())
}

func (gs GrantService) GrantEstimatedExpiry(grantValue *ent.Grant) *time.Time {
	return estimatedExpiry(
		grantValue.ExpiryType,
		gs.config.GrantDefaultExpiryTotalDays,
		gs.config.GrantDefaultExpiryDaysSinceLastUpload,
		grantValue.ExpiryTotalDays,
		grantValue.ExpiryDaysSinceLastUpload,
		grantValue.CreatedAt,
		grantValue.LastUpload,
	)
}

func (gs GrantService) GrantExpired(grantValue *ent.Grant) bool {
	if grantValue.ExpiryType == config.TicketExpiryTypeNone {
		return false
	}
	if grantValue.ExpiryType == config.TicketExpiryTypeSingle {
		return grantValue.TimesUploaded > 0
	}
	estimatedExpiryValue := *gs.GrantEstimatedExpiry(grantValue)
	now := time.Now()

	if grantValue.ExpiryType == config.TicketExpiryTypeCustom {
		return grantValue.TimesUploaded >= uint64(grantValue.ExpiryTotalUploads) ||
			estimatedExpiryValue.Before(now)
	}

	return grantValue.TimesUploaded >= uint64(gs.config.GrantDefaultExpiryTotalUploads) ||
		estimatedExpiryValue.Before(now)

}

func (gs GrantService) ShouldDeleteGrant(grantValue *ent.Grant) bool {
	return len(grantValue.Edges.Files) == 0 && gs.GrantExpired(grantValue)
}

type PublicGrant struct {
	Id              uuid.UUID    `json:"id"`
	Comment         *string      `json:"comment"`
	EstimatedExpiry *string      `json:"estimatedExpiry"`
	User            PublicUser   `json:"owner"`
	Files           []PublicFile `json:"files"`
	CreatedAt       string       `json:"createdAt"`
}

func (gs GrantService) ToPublicGrant(
	fs FileService,
	grantValue *ent.Grant,
	files []*ent.File,
) PublicGrant {
	publicFiles := make([]PublicFile, len(files))
	for i, file := range files {
		publicFiles[i] = fs.ToPublicFile(file)
	}
	var estimatedExpiryValue *string = nil

	if estimatedExpiryResult := gs.GrantEstimatedExpiry(grantValue); estimatedExpiryResult != nil {
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

func NewGrantService(c config.Config) GrantService {
	return GrantService{config: c}
}
