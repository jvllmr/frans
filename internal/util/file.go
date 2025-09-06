package util

import (
	"context"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/grant"
	"github.com/jvllmr/frans/internal/ent/ticket"
	"github.com/jvllmr/frans/internal/ent/user"
)

func UserHasFileAccess(ctx context.Context, userValue *ent.User, fileValue *ent.File) bool {
	return userValue.IsAdmin || fileValue.QueryTickets().
		Where(ticket.HasOwnerWith(user.ID(userValue.ID))).
		CountX(ctx) > 0 || fileValue.QueryGrants().
		Where(grant.HasOwnerWith(user.ID(userValue.ID))).
		CountX(ctx) > 0
}
