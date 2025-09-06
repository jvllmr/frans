package oidc

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/util"
)

func (f *FransOidcProvider) MustGetUser(ctx context.Context, userId uuid.UUID) *ent.User {
	return f.db.User.GetX(ctx, userId)
}

func (f *FransOidcProvider) ProvisionUser(
	ctx context.Context,
	claimsData map[string]any,
) (*ent.User, error) {
	userId, err := uuid.Parse(claimsData["sub"].(string))
	if err != nil {
		return nil, err
	}
	groups := util.InterfaceSliceToStringSlice(claimsData["groups"].([]any))
	isAdmin := slices.Contains(groups, "admin")
	username := claimsData["preferred_username"].(string)
	fullName := claimsData["name"].(string)
	email := claimsData["email"].(string)
	user, err := f.db.User.Get(ctx, userId)
	if err != nil {
		user, err = f.db.User.Create().
			SetGroups(groups).
			SetIsAdmin(isAdmin).
			SetUsername(username).
			SetFullName(fullName).
			SetEmail(email).
			SetID(userId).
			Save(ctx)
		if err != nil {
			slog.Error("Could not create User", "err", err)
			return nil, fmt.Errorf("provision user: %w", err)
		} else {
			slog.Info(fmt.Sprintf("Created user %s", user.Username))
		}

	} else {
		user = f.db.User.UpdateOne(user).
			SetGroups(groups).
			SetIsAdmin(isAdmin).
			SetUsername(username).
			SetFullName(fullName).
			SetEmail(email).
			SaveX(ctx)
		slog.Info(fmt.Sprintf("Updated user %s", username))
	}
	return user, nil
}
