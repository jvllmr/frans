package fransCron

import (
	"context"
	"log/slog"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/util"
)

func GrantsLifecycleTask(configValue config.Config) {
	grants := config.DBClient.Grant.Query().WithFiles().AllX(context.Background())
	deletedCount := 0

	for _, grantValue := range grants {
		if util.ShouldDeleteGrant(configValue, grantValue) {
			err := config.DBClient.Grant.DeleteOne(grantValue).Exec(context.Background())
			if err != nil {
				slog.Error("Could not delete grant", "err", err)
				continue
			}

			deletedCount += 1

		}

	}
	slog.Info("Deleted grants", "count", deletedCount)
}
