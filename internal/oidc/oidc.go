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
	config         config.Config
	OidcConfig     oidc.Config
	db             *ent.Client
}

func (fop *FransOidcProvider) EndSessionEndpoint() string {
	return fop.extraEndpoints.EndSessionEndpoint
}

func NewOIDC(configValue config.Config, db *ent.Client) (*FransOidcProvider, error) {
	slog.Info(
		"Connecting with oidc issuer",
		"issuer",
		configValue.OidcIssuer,
		"client_id",
		configValue.OidcClientID,
		"admin_group",
		configValue.OidcAdminGroup,
	)
	oidcProvider, err := oidc.NewProvider(context.Background(), configValue.OidcIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}
	var extraEndpoints oidcProviderExtraEndpoints
	if err := oidcProvider.Claims(&extraEndpoints); err != nil {
		return nil, fmt.Errorf("failed to find extra endpoints in OIDC Provider: %w", err)
	}
	return &FransOidcProvider{
		config:         configValue,
		Provider:       oidcProvider,
		extraEndpoints: extraEndpoints,
		PKCEManager:    NewPKCEManager(),
		OidcConfig: oidc.Config{
			ClientID: configValue.OidcClientID,
		},
		db: db,
	}, nil
}

func (fop *FransOidcProvider) buildRedirectURL(request *http.Request) string {
	return fmt.Sprintf("%s/api/auth/callback", fop.config.GetBaseURL(request))
}

func (fop *FransOidcProvider) NewOauth2Config(request *http.Request) oauth2.Config {
	endpoint := fop.Endpoint()
	endpoint.AuthStyle = oauth2.AuthStyleInParams
	return oauth2.Config{
		ClientID:     fop.config.OidcClientID,
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
		c.SetCookie(config.AuthOriginCookieName, c.Request.URL.String(), 3_600, "", "", true, true)
		c.Redirect(
			http.StatusTemporaryRedirect,
			oauth2Config.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)),
		)
	} else {
		c.Status(http.StatusUnauthorized)
	}

}

func SetAccessTokenCookie(c *gin.Context, accessToken string) {
	c.SetCookie(config.AccessTokenCookieName, accessToken, 2_592_000, "", "", true, true)
}

func SetIdTokenCookie(c *gin.Context, idToken string) {
	c.SetCookie(config.IdTokenCookieName, idToken, 2_592_000, "", "", true, true)
}
