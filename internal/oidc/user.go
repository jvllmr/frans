package oidc

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/util"
	"golang.org/x/oauth2"
)

func (fop *FransOidcProvider) MustGetUser(ctx context.Context, userId uuid.UUID) *ent.User {
	return fop.db.User.GetX(ctx, userId)
}

func (fop *FransOidcProvider) ProvisionUser(
	ctx context.Context,
	idToken *oidc.IDToken,
	tokenSource *oauth2.TokenSource,
) (*ent.User, error) {
	userInfo, err := fop.UserInfo(ctx, *tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user info: %w", err)
	}
	claimsData := make(map[string]any)
	_ = idToken.Claims(&claimsData)
	_ = userInfo.Claims(&claimsData)
	userId, err := uuid.Parse(claimsData["sub"].(string))
	if err != nil {
		return nil, err
	}

	username := claimsData["preferred_username"].(string)

	groups := make([]string, 0)
	if claimsData["groups"] != nil {
		groups = util.InterfaceSliceToStringSlice(claimsData["groups"].([]any))
	} else {
		slog.WarnContext(ctx, "oidcProvider did not provide groups for user", "username", username)
	}
	isAdmin := slices.Contains(groups, fop.config.OidcAdminGroup)

	fullName := claimsData["name"].(string)
	email := claimsData["email"].(string)
	user, err := fop.db.User.Get(ctx, userId)
	if err != nil {
		user, err = fop.db.User.Create().
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
		user = fop.db.User.UpdateOne(user).
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
