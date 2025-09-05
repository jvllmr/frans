package middleware

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/session"
	"github.com/jvllmr/frans/internal/oidc"

	"golang.org/x/oauth2"
)

func missingAuthResponse(c *gin.Context, oauth2Config oauth2.Config, pkceCache *oidc.PKCECache) {
	if pkceCache != nil {
		state, verifier := pkceCache.CreateChallenge()
		log.Printf("verifier %s", verifier)
		c.SetCookie(config.AuthOriginCookieName, c.Request.URL.String(), 3_600, "", "", true, true)
		c.Redirect(
			http.StatusTemporaryRedirect,
			oauth2Config.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)),
		)
	} else {
		c.Status(http.StatusUnauthorized)
	}

}

func Auth(configValue config.Config, redirect *oidc.PKCECache) gin.HandlerFunc {

	return func(c *gin.Context) {
		var err error
		oauth2Config := oidc.CreateOauth2Config(configValue, c.Request)

		accessTokenCookie, err := c.Request.Cookie(config.AccessTokenCookieName)
		if err != nil {
			slog.Debug("Not authenticated", "err", err)
			missingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}

		idTokenCookie, err := c.Request.Cookie(config.IdTokenCookieName)
		if err != nil {
			slog.Debug("Not authenticated", "err", err)
			missingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}
		session, err := config.DBClient.Session.Query().
			WithUser().
			Where(session.IDToken(idTokenCookie.Value)).
			Only(c.Request.Context())
		if err != nil {
			slog.Debug("Not authenticated", "err", err)
			missingAuthResponse(c, oauth2Config, redirect)
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
		tokenSource := oauth2Config.TokenSource(c.Request.Context(), token)

		newToken, err := tokenSource.Token()
		if err != nil {
			slog.Debug("Not authenticated", "err", err)
			missingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}

		if newToken.Expiry.After(token.Expiry) {
			config.DBClient.Session.UpdateOne(session).
				SetExpire(newToken.Expiry).
				SetRefreshToken(newToken.RefreshToken).
				ExecX(c.Request.Context())
			oidc.SetAccessTokenCookie(c, newToken.AccessToken)
		}

		userinfo, err := oidc.OidcProvider.UserInfo(c.Request.Context(), tokenSource)

		if err != nil {
			slog.Debug("Not authenticated", "err", err)
			missingAuthResponse(c, oauth2Config, redirect)
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
			missingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}

		userId, _ := uuid.Parse(userinfo.Subject)
		user, _ := config.DBClient.User.Get(c.Request.Context(), userId)
		c.Set(config.UserGinContext, user)
	}
}

func GetCurrentUser(ctx *gin.Context) *ent.User {
	return ctx.MustGet(config.UserGinContext).(*ent.User)
}
