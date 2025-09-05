package oidc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"golang.org/x/oauth2"
)

var OidcProvider *oidc.Provider
var OidcProviderExtraEndpoints struct {
	EndSessionEndpoint string `json:"end_session_endpoint"`
}

func NewOidcConfig(configValue config.Config) *oidc.Config {
	return &oidc.Config{ClientID: configValue.OidcClientID}
}

func InitOIDC(configValue config.Config) {
	var err error
	OidcProvider, err = oidc.NewProvider(context.Background(), configValue.OidcIssuer)
	if err != nil {
		slog.Error("Failed to create OIDC provider", "err", err)
		os.Exit(1)
	}
	if err := OidcProvider.Claims(&OidcProviderExtraEndpoints); err != nil {
		slog.Error("Failed to find extra endpoints in OIDC Provider", "err", err)
		os.Exit(1)
	}

}

func buildRedirectURL(configValue config.Config, request *http.Request) string {
	return fmt.Sprintf("%s/api/auth/callback", config.GetBaseURL(configValue, request))
}

func CreateOauth2Config(configValue config.Config, request *http.Request) oauth2.Config {
	endpoint := OidcProvider.Endpoint()
	endpoint.AuthStyle = oauth2.AuthStyleInParams
	return oauth2.Config{
		ClientID:     configValue.OidcClientID,
		ClientSecret: "",
		Endpoint:     endpoint,
		RedirectURL:  buildRedirectURL(configValue, request),
		Scopes:       config.OidcScopes,
	}
}

func SetAccessTokenCookie(c *gin.Context, accessToken string) {
	c.SetCookie(config.AccessTokenCookieName, accessToken, 2_592_000, "", "", true, true)
}

func SetIdTokenCookie(c *gin.Context, idToken string) {
	c.SetCookie(config.IdTokenCookieName, idToken, 2_592_000, "", "", true, true)
}
