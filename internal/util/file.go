package util

import (
	"context"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/user"
)

func UserHasFileAccess(ctx context.Context, userValue *ent.User, fileValue *ent.File) bool {
	return userValue.IsAdmin || fileValue.QueryData().QueryUsers().
		Where(user.ID(userValue.ID)).
		CountX(ctx) > 0
}
