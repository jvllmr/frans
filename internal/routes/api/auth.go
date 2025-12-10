package apiRoutes

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/oidc"
	"github.com/jvllmr/frans/internal/otel"
	"github.com/jvllmr/frans/internal/util"
	"golang.org/x/oauth2"
)

type authController struct {
	cfg      config.Config
	provider *oidc.FransOidcProvider
}

func (ac *authController) redirectToAuth(c *gin.Context, oauth2Config oauth2.Config) {
	state, verifier := ac.provider.CreateChallenge(c)

	c.Redirect(
		http.StatusTemporaryRedirect,
		oauth2Config.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)),
	)
}

func (ac *authController) authCallback(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "authCallback")
	defer span.End()
	code := c.Request.URL.Query().Get("code")
	oauth2Config := ac.provider.NewOauth2Config(c.Request)

	pkceVerifier, err := ac.provider.GetVerifier(c)

	if err != nil {
		slog.Error("Failed at Oauth2 token exchange", "err", err)
		ac.redirectToAuth(c, oauth2Config)
		return
	}

	oauth2Token, err := oauth2Config.Exchange(
		ctx,
		code,
		oauth2.VerifierOption(pkceVerifier),
	)
	if err != nil {
		slog.Error("Failed at OAuth2 token exchange", "err", err)
		ac.redirectToAuth(c, oauth2Config)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		slog.Error("Failed to retrieve id_token from access token", "err", err)
		ac.redirectToAuth(c, oauth2Config)
		return
	}

	verifier := ac.provider.Verifier(&ac.provider.OidcConfig)
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		slog.Error("Failed to verify id token", "err", err)
		ac.redirectToAuth(c, oauth2Config)
		return
	}
	tokenSource := oauth2Config.TokenSource(ctx, oauth2Token)

	user, err := ac.provider.ProvisionUser(ctx, idToken, &tokenSource)
	if err != nil {
		slog.Error("Could not provision user", "err", err)
		util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
		return
	}

	ac.provider.SetIdTokenCookie(c, rawIDToken)
	ac.provider.SetAccessTokenCookie(c, oauth2Token.AccessToken)
	err = ac.provider.CreateSession(ctx, user, oauth2Token, rawIDToken)
	if err != nil {
		slog.Error("could not store session", "err", err)
	}
	authOrigin, err := c.Request.Cookie(config.AuthOriginCookieName)

	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, ac.cfg.RootPath)
		return
	}
	authOriginPath, err := url.PathUnescape(authOrigin.Value)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, ac.cfg.RootPath)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, authOriginPath)
}

func (ac *authController) logoutCallback(c *gin.Context) {
	c.SetCookie(config.AccessTokenCookieName, "", 5, ac.cfg.RootPath, "", true, true)
	idTokenCookie, err := c.Request.Cookie(config.IdTokenCookieName)
	if err == nil {
		_ = ac.provider.DeleteSession(c.Request.Context(), idTokenCookie)
	}
	c.SetCookie(config.IdTokenCookieName, "", 5, ac.cfg.RootPath, "", true, true)
	redirectURI := c.Query("redirect_uri")
	slog.Debug(fmt.Sprintf("logout redirect %s", redirectURI))
	c.Redirect(
		http.StatusTemporaryRedirect,
		fmt.Sprintf(
			"%s?id_token_hint=%s&post_logout_redirect_uri=%s",
			ac.provider.EndSessionEndpoint(),
			idTokenCookie.Value,
			redirectURI,
		),
	)
}

func setupAuthRoutes(
	r *gin.RouterGroup,
	configValue config.Config,

	provider *oidc.FransOidcProvider,
) {
	authGroup := r.Group("/auth")
	controller := authController{
		cfg:      configValue,
		provider: provider,
	}

	authGroup.GET("/callback", controller.authCallback)

	authGroup.GET("/logout", controller.logoutCallback)

}
