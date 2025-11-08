package tasks

import (
	"testing"
	"time"

	"github.com/jvllmr/frans/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestSessionLifecycleTask(t *testing.T) {
	db := testutil.SetupTestDBClient(t)
	now := time.Now()
	shareAccessToken := db.ShareAccessToken.Create().
		SetID("not_expired").
		SetExpiry(now.Add(time.Hour)).
		SaveX(t.Context())
	_ = db.ShareAccessToken.Create().
		SetID("expired").
		SetExpiry(now.Add(-time.Minute)).
		SaveX(t.Context())
	session := db.Session.Create().
		SetIDToken("session1").
		SetRefreshToken("dummy_refresh").
		SetExpire(now).
		SaveX(t.Context())
	_ = db.Session.Create().
		SetIDToken("session2").
		SetRefreshToken("dummy_refresh").
		SetExpire(now.Add(-time.Hour)).
		SaveX(t.Context())
	SessionLifecycleTask(db)

	remainingTokens := db.ShareAccessToken.Query().AllX(t.Context())
	assert.Equal(t, 1, len(remainingTokens))
	assert.Equal(t, shareAccessToken.ID, remainingTokens[0].ID)

	remainingSessions := db.Session.Query().AllX(t.Context())
	assert.Equal(t, 1, len(remainingSessions))
	assert.Equal(t, session.ID, remainingSessions[0].ID)
}
