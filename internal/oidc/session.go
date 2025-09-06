package oidc

import (
	"context"
	"net/http"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/session"
	"golang.org/x/oauth2"
)

func (f *FransOidcProvider) GetSession(
	ctx context.Context,
	idTokenCookie *http.Cookie,
) (*ent.Session, error) {
	return f.db.Session.Query().
		WithUser().
		Where(session.IDToken(idTokenCookie.Value)).
		Only(ctx)
}

func (f *FransOidcProvider) UpdateSession(
	ctx context.Context,
	session *ent.Session,
	newToken *oauth2.Token,
) {
	f.db.Session.UpdateOne(session).
		SetExpire(newToken.Expiry).
		SetRefreshToken(newToken.RefreshToken).
		ExecX(ctx)
}

func (f *FransOidcProvider) DeleteSession(ctx context.Context, idTokenCookie *http.Cookie) error {
	_, err := f.db.Session.Delete().
		Where(session.IDToken(idTokenCookie.Value)).
		Exec(ctx)
	return err
}

func (f *FransOidcProvider) CreateSession(
	ctx context.Context,
	user *ent.User,
	token *oauth2.Token,
	rawIDToken string,
) error {
	return f.db.Session.Create().
		SetUser(user).
		SetExpire(token.Expiry).
		SetIDToken(rawIDToken).
		SetRefreshToken(token.RefreshToken).
		Exec(ctx)
}
