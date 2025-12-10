package oidc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"golang.org/x/oauth2"
)

type oidcProviderExtraEndpoints struct {
	EndSessionEndpoint string `json:"end_session_endpoint"`
}

type FransOidcProvider struct {
	*oidc.Provider
	*PKCEManager
	extraEndpoints oidcProviderExtraEndpoints
	cfg            config.Config
	OidcConfig     oidc.Config
	db             *ent.Client
}

func (fop *FransOidcProvider) EndSessionEndpoint() string {
	return fop.extraEndpoints.EndSessionEndpoint
}

func NewOIDC(cfg config.Config, db *ent.Client) (*FransOidcProvider, error) {
	slog.Info(
		"Connecting with oidc issuer",
		"issuer",
		cfg.OidcIssuer,
		"client_id",
		cfg.OidcClientID,
		"admin_group",
		cfg.OidcAdminGroup,
	)
	oidcProvider, err := oidc.NewProvider(context.Background(), cfg.OidcIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}
	var extraEndpoints oidcProviderExtraEndpoints
	if err := oidcProvider.Claims(&extraEndpoints); err != nil {
		return nil, fmt.Errorf("failed to find extra endpoints in OIDC Provider: %w", err)
	}
	return &FransOidcProvider{
		cfg:            cfg,
		Provider:       oidcProvider,
		extraEndpoints: extraEndpoints,
		PKCEManager:    NewPKCEManager(cfg),
		OidcConfig: oidc.Config{
			ClientID: cfg.OidcClientID,
		},
		db: db,
	}, nil
}

func (fop *FransOidcProvider) buildRedirectURL(request *http.Request) string {
	return fmt.Sprintf("%s/api/auth/callback", fop.cfg.GetBaseURL(request))
}

func (fop *FransOidcProvider) NewOauth2Config(request *http.Request) oauth2.Config {
	endpoint := fop.Endpoint()
	endpoint.AuthStyle = oauth2.AuthStyleInParams
	return oauth2.Config{
		ClientID:     fop.cfg.OidcClientID,
		ClientSecret: "",
		Endpoint:     endpoint,
		RedirectURL:  fop.buildRedirectURL(request),
		Scopes:       config.OidcScopes,
	}
}

func (fop *FransOidcProvider) MissingAuthResponse(
	c *gin.Context,
	oauth2Config oauth2.Config,
	redirect bool,
) {
	if redirect {
		state, verifier := fop.CreateChallenge(c)
		c.SetCookie(
			config.AuthOriginCookieName,
			c.Request.URL.String(),
			3_600,
			fop.cfg.RootPath,
			"",
			true,
			true,
		)
		c.Redirect(
			http.StatusTemporaryRedirect,
			oauth2Config.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)),
		)
	} else {
		c.Status(http.StatusUnauthorized)
	}

}

func (fop *FransOidcProvider) SetAccessTokenCookie(c *gin.Context, accessToken string) {
	c.SetCookie(
		config.AccessTokenCookieName,
		accessToken,
		2_592_000,
		fop.cfg.RootPath,
		"",
		true,
		true,
	)
}

func (fop *FransOidcProvider) SetIdTokenCookie(c *gin.Context, idToken string) {
	c.SetCookie(config.IdTokenCookieName, idToken, 2_592_000, fop.cfg.RootPath, "", true, true)
}
