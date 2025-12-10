package oidc

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"golang.org/x/oauth2"
)

type OidcState = string
type oidcVerifier = string

type PKCEManager struct {
	cfg config.Config
}

func (p *PKCEManager) CreateChallenge(c *gin.Context) (OidcState, oidcVerifier) {
	var state = uuid.New().String()
	verifier := oauth2.GenerateVerifier()
	c.SetCookie(state, verifier, 3_600, p.cfg.RootPath, "", true, true)
	return state, verifier
}

func (p *PKCEManager) GetVerifier(c *gin.Context) (oidcVerifier, error) {
	var state = c.Query("state")
	verifier, err := c.Cookie(state)
	if err != nil {
		return "", fmt.Errorf("could not retrieve PKCE verifier for given state")
	}
	c.SetCookie(state, "", 0, p.cfg.RootPath, "", true, true)
	return verifier, nil
}

func NewPKCEManager(cfg config.Config) *PKCEManager {
	return &PKCEManager{cfg: cfg}
}
