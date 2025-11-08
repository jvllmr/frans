package apiRoutes

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/util"
)

func createTestGrant(
	t *testing.T,
	db *ent.Client,
	user *ent.User,
	modifier func(q *ent.GrantCreate) *ent.GrantCreate,
) *ent.Grant {
	if modifier == nil {
		modifier = func(q *ent.GrantCreate) *ent.GrantCreate { return q }
	}
	salt := util.GenerateSalt()
	hashedPassword := util.HashPassword("abc123", salt)
	grantValue := modifier(
		db.Grant.Create().
			SetID(uuid.New()).
			SetComment("").
			SetCreatorLang("en").
			SetEmailOnUpload("testmail@vllmr.dev").
			SetExpiryDaysSinceLastUpload(7).
			SetExpiryTotalDays(30).
			SetExpiryTotalUploads(10).
			SetExpiryType("auto").
			SetFileExpiryDaysSinceLastDownload(7).
			SetFileExpiryTotalDays(30).
			SetFileExpiryTotalDownloads(10).
			SetFileExpiryType("auto").
			SetHashedPassword(hashedPassword).SetOwner(user).SetSalt(string(salt)),
	).SaveX(t.Context())

	return grantValue
}
