package config

import (
	"github.com/coreos/go-oidc/v3/oidc"
)

var OidcScopes = []string{oidc.ScopeOpenID, "profile", "email"}

const AccessTokenCookieName = "frans_access_token"
const IdTokenCookieName = "frans_id_token"
const AuthOriginCookieName = "frans_auth_origin"

const UserGinContext = "user"
