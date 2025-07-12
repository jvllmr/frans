package config

import (
	"github.com/coreos/go-oidc/v3/oidc"
)

var OidcScopes = []string{oidc.ScopeOpenID, "profile", "email"}

const AccessTokenCookieName = "frans_access_token"
const IdTokenCookieName = "frans_id_token"
const AuthOriginCookieName = "frans_auth_origin"
const ShareAccessTokenCookieName = "frans_share_access_token"

const UserGinContext = "user"
const ShareTicketContext = "shareTicket"

const TicketExpiryTypeAuto = "auto"
const TicketExpiryTypeSingle = "single"
const TicketExpiryTypeNone = "none"
const TicketExpiryTypeCustom = "custom"

const ShareAccessTokenExpirySeconds = 10
