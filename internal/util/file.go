package util

import (
	"context"

	"codeberg.org/jvllmr/frans/internal/ent"
	"codeberg.org/jvllmr/frans/internal/ent/user"
)

func UserHasFileAccess(ctx context.Context, userValue *ent.User, fileValue *ent.File) bool {
	return userValue.IsAdmin || fileValue.QueryOwner().
		Where(user.ID(userValue.ID)).
		CountX(ctx) > 0
}
