package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var OidcProvider *oidc.Provider
var OidcProviderExtraEndpoints struct {
	EndSessionEndpoint    string `json:"end_session_endpoint"`
	IntrospectionEndpoint string `json:"introspection_endpoint"`
}

func doOidcRequest(
	ctx context.Context,
	configValue Config,
	url string,
	data url.Values,
) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(configValue.OidcClientID, configValue.OidcClientSecret)
	return http.DefaultClient.Do(request)
}

type IntrospectionResponse struct {
	Active bool   `json:"active"`
	Exp    int64  `json:"exp"`
	Sub    string `json:"sub"`
}

func DoIntrospection(
	ctx context.Context,
	configValue Config,
	accessToken string,
) (*IntrospectionResponse, error) {
	data := url.Values{}
	data.Set("token", accessToken)
	response, err := doOidcRequest(
		ctx,
		configValue,
		OidcProviderExtraEndpoints.IntrospectionEndpoint,
		data,
	)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf(
			"introspection response returned with status code %s",
			response.Status,
		)
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result *IntrospectionResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type TokenRefreshmentResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func DoTokenRefreshment(
	ctx context.Context,
	configValue Config,
	refreshToken string,
) (*TokenRefreshmentResponse, error) {
	data := url.Values{}
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", configValue.OidcClientID)
	data.Set("client_secret", configValue.OidcClientSecret)
	data.Set("scope", strings.Join(OidcScopes, " "))
	response, err := doOidcRequest(ctx, configValue, OidcProvider.Endpoint().TokenURL, data)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf(
			"token refreshment response returned with status code %s",
			response.Status,
		)
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result *TokenRefreshmentResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func InitOIDC(config Config) {
	var err error
	OidcProvider, err = oidc.NewProvider(context.Background(), config.OidcIssuer)
	if err != nil {
		log.Fatalf("Failed to create OIDC provider: %v", err)
	}
	if err := OidcProvider.Claims(&OidcProviderExtraEndpoints); err != nil {
		log.Fatalf("Failed to find extra endpoints in OIDC Provider: %v", err)
	}

}

func buildRedirectURL(config Config, request *http.Request) string {
	proto := "http"
	if request.TLS != nil {
		proto = "https"
	}
	host := request.Host
	patchedRootPath := config.RootPath
	if len(patchedRootPath) == 0 {
		patchedRootPath = "/"
	}
	return fmt.Sprintf("%s://%s%s/api/auth/callback", proto, host, patchedRootPath)
}

func CreateOauth2Config(configValue Config, request *http.Request) oauth2.Config {
	return oauth2.Config{
		ClientID:     configValue.OidcClientID,
		ClientSecret: configValue.OidcClientSecret,
		Endpoint:     OidcProvider.Endpoint(),
		RedirectURL:  buildRedirectURL(configValue, request),
		Scopes:       OidcScopes,
	}
}

func SetAccessTokenCookie(c *gin.Context, accessToken string) {
	c.SetCookie(AccessTokenCookieName, accessToken, 2_592_000, "", "", true, true)
}

func SetIdTokenCookie(c *gin.Context, idToken string) {
	c.SetCookie(IdTokenCookieName, idToken, 2_592_000, "", "", true, true)
}
