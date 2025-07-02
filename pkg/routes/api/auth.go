package apiRoutes

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent/session"
	"github.com/jvllmr/frans/pkg/util"
)

func setupAuthRoutes(r *gin.RouterGroup, configValue config.Config) {
	authGroup := r.Group("/auth")
	// @Router /api/v1/auth/callback
	authCallback := func(c *gin.Context) {
		code := c.Request.URL.Query().Get("code")
		oauth2Config := config.CreateOauth2Config(configValue, c.Request)
		oauth2Token, err := oauth2Config.Exchange(c.Request.Context(), code)
		if err != nil {
			log.Printf("Error :Failed at OAuth2 token exchange: %v", err)
			c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(""))
			return
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			log.Printf("Error: Failed to retrieve id_token from access token: %s", err)
			c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(""))
			return
		}

		verifier := config.OidcProvider.Verifier(&oidc.Config{ClientID: configValue.OidcClientID})
		idToken, err := verifier.Verify(c.Request.Context(), rawIDToken)
		if err != nil {
			log.Printf("Error: Failed to verify id token: %v", err)
			c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(""))
			return
		}
		tokenSource := oauth2Config.TokenSource(c.Request.Context(), oauth2Token)
		userInfo, err := config.OidcProvider.UserInfo(c.Request.Context(), tokenSource)
		if err != nil {
			log.Printf("Error: Failed to retrieve user info: %v", err)
			c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(""))
			return
		}
		claimsData := make(map[string]any)
		_ = idToken.Claims(&claimsData)
		_ = userInfo.Claims(&claimsData)

		userId, err := uuid.Parse(claimsData["sub"].(string))
		if err != nil {
			log.Printf("Error: Could not access user id from claims: %v", err)
			c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(""))
			return
		}
		groups := util.InterfaceSliceToStringSlice(claimsData["groups"].([]interface{}))
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
				log.Printf("Error: Could not create User %s", err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			} else {
				log.Printf("Info: Created user %s", user.Username)
			}

		} else {
			_ = config.DBClient.User.UpdateOne(user).SetGroups(groups).SetIsAdmin(isAdmin).SetUsername(username).SetFullName(fullName).SetEmail(email).Exec(c.Request.Context())
			log.Printf("Info: Updated user %s", username)
		}
		config.SetIdTokenCookie(c, rawIDToken)
		config.SetAccessTokenCookie(c, oauth2Token.AccessToken)
		err = config.DBClient.Session.Create().
			SetUser(user).
			SetExpire(oauth2Token.Expiry).
			SetIDToken(rawIDToken).
			SetRefreshToken(oauth2Token.RefreshToken).
			Exec(c.Request.Context())
		if err != nil {
			log.Printf("Error: could not store session: %v", err)
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
		log.Printf("redirect %s", redirectURI)
		ctx.Redirect(
			http.StatusTemporaryRedirect,
			fmt.Sprintf(
				"%s?id_token_hint=%s&post_logout_redirect_uri=%s",
				config.OidcProviderExtraEndpoints.EndSessionEndpoint,
				idTokenCookie.Value,
				redirectURI,
			),
		)
	})

}
