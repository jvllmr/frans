package routesUtil

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent/session"
	"golang.org/x/oauth2"
)

func missingAuthResponse(c *gin.Context, oauth2Config oauth2.Config, redirect bool) {
	if redirect {
		c.SetCookie(config.AuthOriginCookieName, c.Request.URL.String(), 3_600, "", "", true, true)
		c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(""))
	} else {
		c.Status(http.StatusUnauthorized)
	}

}

func AuthMiddleware(configValue config.Config, redirect bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		_ = config.DBClient.Session.Delete().Where(session.ExpireLT(time.Now().Add(-1 * time.Hour)))
		var err error
		oauth2Config := config.CreateOauth2Config(configValue, c.Request)

		accessTokenCookie, err := c.Request.Cookie(config.AccessTokenCookieName)
		if err != nil {
			log.Printf("Not authenticated: %v", err)
			missingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}
		introspectionResponse, err := config.DoIntrospection(
			c.Request.Context(),
			configValue,
			accessTokenCookie.Value,
		)
		if err != nil {
			log.Printf("Not authenticated: %v", err)
			missingAuthResponse(c, oauth2Config, redirect)
			c.Abort()
			return
		}

		if !introspectionResponse.Active {
			idTokenCookie, err := c.Request.Cookie(config.IdTokenCookieName)
			if err != nil {
				log.Printf("Not authenticated: %v", err)
				missingAuthResponse(c, oauth2Config, redirect)
				c.Abort()
				return
			}
			session, err := config.DBClient.Session.Query().
				WithUser().
				Where(session.IDToken(idTokenCookie.Value)).
				Only(c.Request.Context())
			if err != nil {
				log.Printf("Not authenticated: %v", err)
				missingAuthResponse(c, oauth2Config, redirect)
				c.Abort()
				return
			}
			tokenRefreshmentResponse, err := config.DoTokenRefreshment(
				c.Request.Context(),
				configValue,
				session.RefreshToken,
			)
			if err != nil {
				log.Printf("Not authenticated: %v", err)
				missingAuthResponse(c, oauth2Config, redirect)
				c.Abort()
				return
			}
			introspectionResponse, err = config.DoIntrospection(
				c.Request.Context(),
				configValue,
				tokenRefreshmentResponse.AccessToken,
			)
			if err != nil {
				log.Printf("Not authenticated: %v", err)
				missingAuthResponse(c, oauth2Config, redirect)
				c.Abort()
				return
			}
			if introspectionResponse.Sub != session.Edges.User.ID.String() {
				log.Printf(
					"Not authenticated: Introspection sub did not match user id - %s != %s",
					introspectionResponse.Sub,
					session.Edges.User.ID,
				)
				missingAuthResponse(c, oauth2Config, redirect)
				c.Abort()
				return
			}
			config.SetAccessTokenCookie(c, tokenRefreshmentResponse.AccessToken)
			newExpiry := time.Now().
				Add(time.Duration(tokenRefreshmentResponse.ExpiresIn) * time.Second)
			_ = config.DBClient.Session.UpdateOne(session).
				SetExpire(newExpiry).
				SetRefreshToken(tokenRefreshmentResponse.RefreshToken).
				Exec(c.Request.Context())
		}
		userId, _ := uuid.Parse(introspectionResponse.Sub)
		user, _ := config.DBClient.User.Get(c.Request.Context(), userId)
		c.Set(config.UserGinContext, user)
	}
}
