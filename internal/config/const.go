package config

import (
	"github.com/coreos/go-oidc/v3/oidc"
)

var OidcScopes = []string{oidc.ScopeOpenID, "profile", "email"}

const (
	AccessTokenCookieName      = "frans_access_token"
	IdTokenCookieName          = "frans_id_token"
	AuthOriginCookieName       = "frans_auth_origin"
	ShareAccessTokenCookieName = "frans_share_access_token"
)

const (
	UserGinContext     = "user"
	ShareTicketContext = "shareTicket"
	ShareGrantContext  = "shareGrant"
)

const (
	TicketExpiryTypeAuto   = "auto"
	TicketExpiryTypeSingle = "single"
	TicketExpiryTypeNone   = "none"
	TicketExpiryTypeCustom = "custom"
)

const ShareAccessTokenExpirySeconds = 10
