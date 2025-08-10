package fransCron

import (
	"context"
	"log/slog"
	"time"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent/session"
	"github.com/jvllmr/frans/internal/ent/shareaccesstoken"
)

func SessionLifecycleTask() {
	now := time.Now()

	deletedTokens := config.DBClient.ShareAccessToken.Delete().
		Where(shareaccesstoken.ExpiryLT(now)).
		ExecX(context.Background())
	slog.Info("Deleted shared access tokens", "count", deletedTokens)
	deletedSessions := config.DBClient.Session.Delete().
		Where(session.ExpireLT(now)).
		ExecX(context.Background())
	slog.Info("Deleted sessions", "count", deletedSessions)
}
