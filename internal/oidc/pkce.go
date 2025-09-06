package oidc

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"golang.org/x/oauth2"
)

type OidcState = string
type oidcVerifier = string

type PKCECache struct {
	lru *expirable.LRU[OidcState, oidcVerifier]
}

func (p *PKCECache) CreateChallenge() (OidcState, oidcVerifier) {
	var state OidcState = uuid.New().String()
	verifier := oauth2.GenerateVerifier()

	p.lru.Add(state, verifier)

	return state, verifier
}

func (p *PKCECache) GetVerifier(state OidcState) (oidcVerifier, error) {
	verifier, ok := p.lru.Get(state)
	if !ok {
		return "", fmt.Errorf("frans does not know PKCE verifier for given state")
	}
	return verifier, nil
}

func NewPKCECache() *PKCECache {
	return &PKCECache{lru: expirable.NewLRU[string, string](2048, nil, 1*time.Hour)}
}
