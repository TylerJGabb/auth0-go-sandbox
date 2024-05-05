package authenticator

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
	JwtKeyGetter
}

// New instantiates the *Authenticator.
func New() (*Authenticator, error) {
	domain := os.Getenv("AUTH0_DOMAIN")
	issuer := fmt.Sprintf("https://%s/", domain)
	wellKnownUrl := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)
	provider, err := oidc.NewProvider(
		context.Background(),
		issuer,
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	keyGetter := NewKeyGetter(wellKnownUrl)
	return &Authenticator{
		Provider:     provider,
		Config:       conf,
		JwtKeyGetter: keyGetter,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

func (a *Authenticator) VerifyToken(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.Parse(token, a.JwtKeyGetter.GetKey)
}
