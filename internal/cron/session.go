package fransCron

import (
	"context"
	"log/slog"
	"time"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/session"
	"github.com/jvllmr/frans/internal/ent/shareaccesstoken"
)

func SessionLifecycleTask(db *ent.Client) {
	now := time.Now()

	deletedTokens := db.ShareAccessToken.Delete().
		Where(shareaccesstoken.ExpiryLT(now)).
		ExecX(context.Background())
	slog.Info("Deleted shared access tokens", "count", deletedTokens)
	deletedSessions := db.Session.Delete().
		Where(session.ExpireLT(now.Add(-1 * time.Hour))).
		ExecX(context.Background())
	slog.Info("Deleted sessions", "count", deletedSessions)
}
