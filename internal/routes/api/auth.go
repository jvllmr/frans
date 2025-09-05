package apiRoutes

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent/session"
	"github.com/jvllmr/frans/internal/oidc"
	"github.com/jvllmr/frans/internal/util"
	"golang.org/x/oauth2"
)

func redirectToAuth(c *gin.Context, oauth2Config oauth2.Config, pkceCache *oidc.PKCECache) {
	state, verifier := pkceCache.CreateChallenge()

	c.Redirect(
		http.StatusTemporaryRedirect,
		oauth2Config.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)),
	)
}

func setupAuthRoutes(r *gin.RouterGroup, configValue config.Config, pkceCache *oidc.PKCECache) {
	authGroup := r.Group("/auth")
	// @Router /api/v1/auth/callback
	authCallback := func(c *gin.Context) {
		code := c.Request.URL.Query().Get("code")
		oauth2Config := oidc.CreateOauth2Config(configValue, c.Request)

		var state oidc.OidcState = c.Request.URL.Query().Get("state")
		pkceVerifier, err := pkceCache.GetVerifier(state)

		if err != nil {
			slog.Error("Failed at Oauth2 token exchange", "err", err)
			redirectToAuth(c, oauth2Config, pkceCache)
			return
		}

		oauth2Token, err := oauth2Config.Exchange(
			c.Request.Context(),
			code,
			oauth2.VerifierOption(pkceVerifier),
		)
		if err != nil {
			slog.Error("Failed at OAuth2 token exchange", "err", err)
			redirectToAuth(c, oauth2Config, pkceCache)
			return
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			slog.Error("Failed to retrieve id_token from access token", "err", err)
			redirectToAuth(c, oauth2Config, pkceCache)
			return
		}

		verifier := oidc.OidcProvider.Verifier(oidc.NewOidcConfig(configValue))
		idToken, err := verifier.Verify(c.Request.Context(), rawIDToken)
		if err != nil {
			slog.Error("Failed to verify id token", "err", err)
			redirectToAuth(c, oauth2Config, pkceCache)
			return
		}
		tokenSource := oauth2Config.TokenSource(c.Request.Context(), oauth2Token)

		userInfo, err := oidc.OidcProvider.UserInfo(c.Request.Context(), tokenSource)
		if err != nil {
			slog.Error("Failed to retrieve user info", "err", err)
			redirectToAuth(c, oauth2Config, pkceCache)
			return
		}
		claimsData := make(map[string]any)
		_ = idToken.Claims(&claimsData)
		_ = userInfo.Claims(&claimsData)

		userId, err := uuid.Parse(claimsData["sub"].(string))
		if err != nil {
			slog.Error("Could not access user id from claims", "err", err)
			redirectToAuth(c, oauth2Config, pkceCache)
			return
		}
		groups := util.InterfaceSliceToStringSlice(claimsData["groups"].([]any))
		isAdmin := slices.Contains(groups, "admin")
		username := claimsData["preferred_username"].(string)
		fullName := claimsData["name"].(string)
		email := claimsData["email"].(string)
		user, err := config.DBClient.User.Get(c.Request.Context(), userId)
		if err != nil {
			user, err = config.DBClient.User.Create().
				SetGroups(groups).
				SetIsAdmin(isAdmin).
				SetUsername(username).
				SetFullName(fullName).
				SetEmail(email).
				SetID(userId).
				Save(c.Request.Context())
			if err != nil {
				slog.Error("Could not create User", "err", err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			} else {
				slog.Info(fmt.Sprintf("Created user %s", user.Username))
			}

		} else {
			_ = config.DBClient.User.UpdateOne(user).
				SetGroups(groups).
				SetIsAdmin(isAdmin).
				SetUsername(username).
				SetFullName(fullName).
				SetEmail(email).
				Exec(c.Request.Context())
			slog.Info(fmt.Sprintf("Updated user %s", username))
		}
		oidc.SetIdTokenCookie(c, rawIDToken)
		oidc.SetAccessTokenCookie(c, oauth2Token.AccessToken)
		err = config.DBClient.Session.Create().
			SetUser(user).
			SetExpire(oauth2Token.Expiry).
			SetIDToken(rawIDToken).
			SetRefreshToken(oauth2Token.RefreshToken).
			Exec(c.Request.Context())
		if err != nil {
			slog.Error("could not store session", "err", err)
		}
		authOrigin, err := c.Request.Cookie(config.AuthOriginCookieName)

		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, configValue.RootPath)
			return
		}
		authOriginPath, err := url.PathUnescape(authOrigin.Value)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, configValue.RootPath)
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, authOriginPath)

	}

	authGroup.GET("/callback", authCallback)

	authGroup.GET("/logout", func(ctx *gin.Context) {
		ctx.SetCookie(config.AccessTokenCookieName, "", 5, "", "", true, true)
		idTokenCookie, err := ctx.Request.Cookie(config.IdTokenCookieName)
		if err == nil {
			_, _ = config.DBClient.Session.Delete().
				Where(session.IDToken(idTokenCookie.Value)).
				Exec(ctx.Request.Context())
		}
		ctx.SetCookie(config.IdTokenCookieName, "", 5, "", "", true, true)
		redirectURI := ctx.Query("redirect_uri")
		slog.Debug(fmt.Sprintf("logout redirect %s", redirectURI))
		ctx.Redirect(
			http.StatusTemporaryRedirect,
			fmt.Sprintf(
				"%s?id_token_hint=%s&post_logout_redirect_uri=%s",
				oidc.OidcProviderExtraEndpoints.EndSessionEndpoint,
				idTokenCookie.Value,
				redirectURI,
			),
		)
	})

}
