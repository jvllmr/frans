package oidc

import (
	"context"
	"net/http"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/session"
	"golang.org/x/oauth2"
)

func (fop *FransOidcProvider) GetSession(
	ctx context.Context,
	idTokenCookie *http.Cookie,
) (*ent.Session, error) {
	return fop.db.Session.Query().
		WithUser().
		Where(session.IDToken(idTokenCookie.Value)).
		Only(ctx)
}

func (fop *FransOidcProvider) UpdateSession(
	ctx context.Context,
	session *ent.Session,
	newToken *oauth2.Token,
) {
	fop.db.Session.UpdateOne(session).
		SetExpire(newToken.Expiry).
		SetRefreshToken(newToken.RefreshToken).
		ExecX(ctx)
}

func (fop *FransOidcProvider) DeleteSession(ctx context.Context, idTokenCookie *http.Cookie) error {
	_, err := fop.db.Session.Delete().
		Where(session.IDToken(idTokenCookie.Value)).
		Exec(ctx)
	return err
}

func (fop *FransOidcProvider) CreateSession(
	ctx context.Context,
	user *ent.User,
	token *oauth2.Token,
	rawIDToken string,
) error {
	return fop.db.Session.Create().
		SetUser(user).
		SetExpire(token.Expiry).
		SetIDToken(rawIDToken).
		SetRefreshToken(token.RefreshToken).
		Exec(ctx)
}
