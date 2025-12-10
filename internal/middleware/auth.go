package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/oidc"
	"github.com/jvllmr/frans/internal/otel"

	"golang.org/x/oauth2"
)

func Auth(p *oidc.FransOidcProvider, redirect bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, span := otel.NewSpan(c.Request.Context(), "checkAuth")
		defer span.End()
		var err error
		oauth2Config := p.NewOauth2Config(c.Request)

		accessTokenCookie, err := c.Request.Cookie(config.AccessTokenCookieName)
		if err != nil {
			slog.Debug("Not authenticated", "err", err)
			p.MissingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}

		idTokenCookie, err := c.Request.Cookie(config.IdTokenCookieName)
		if err != nil {
			slog.Debug("Not authenticated", "err", err)
			p.MissingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}
		session, err := p.GetSession(ctx, idTokenCookie)
		if err != nil {
			slog.Debug("Not authenticated", "err", err)
			p.MissingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}
		now := time.Now()
		token := &oauth2.Token{
			AccessToken:  accessTokenCookie.Value,
			TokenType:    "bearer",
			RefreshToken: session.RefreshToken,
			Expiry:       session.Expire,
			ExpiresIn:    session.Expire.Unix() - now.Unix(),
		}
		tokenSource := oauth2Config.TokenSource(ctx, token)

		newToken, err := tokenSource.Token()
		if err != nil {
			slog.Debug("Not authenticated", "err", err)
			p.MissingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}

		if newToken.Expiry.After(token.Expiry) {
			p.UpdateSession(ctx, session, newToken)
			p.SetAccessTokenCookie(c, newToken.AccessToken)
		}

		userinfo, err := p.UserInfo(ctx, tokenSource)

		if err != nil {
			slog.Debug("Not authenticated", "err", err)
			p.MissingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}

		if userinfo.Subject != session.Edges.User.ID.String() {
			slog.Warn(
				"Not authenticated: Userinfo sub did not match user id",
				"userinfo_sub",
				userinfo.Subject,
				"user_id",
				session.Edges.User.ID,
			)
			p.MissingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}

		userId, _ := uuid.Parse(userinfo.Subject)
		user := p.MustGetUser(ctx, userId)
		c.Set(config.UserGinContext, user)
	}
}

func GetCurrentUser(ctx *gin.Context) *ent.User {
	return ctx.MustGet(config.UserGinContext).(*ent.User)
}
