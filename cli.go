package byoidccli

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Config struct {
	OidcProviderURL string
	ClientID        string
	ClientSecret    string
	RedirectURL     string
}

type ResponseTokens struct {
	AccessToken  string `json:"access_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`

	Claims IDTokenClaims `json:"claims,omitempty"`

	TokenType string `json:"token_type,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
	Scope     string `json:"scope,omitempty"`
}

type IDTokenClaims struct {
	*BasicClaims

	Audience  string `json:"aud,omitempty"`     // ClientID
	Nounce    string `json:"nounce,omitempty"`  // Client 生成的随机值
	ATHash    string `json:"at_hash,omitempty"` // Access Token Hash, 用于Client验证AccessToken
	Subject   string `json:"sub,omitempty"`     // UserID, 用户唯一标识
	Role      string `json:"rol,omitempty"`     // Role of User
	SessionId string `json:"sid,omitempty"`

	*UserinfoClaims
}

type BasicClaims struct {
	Issuer    string `json:"iss,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"` // Token发放时间。以 Unix 时间戳表示。
	ExpiresAt int64  `json:"exp,omitempty"` // Token过期时间。以 Unix 时间戳表示。
}

type UserinfoClaims struct {
	Nickname string `json:"nickname,omitempty"`
	Picture  string `json:"picture,omitempty"`
	Gender   string `json:"gender,omitempty"`

	SchoolNumber  string `json:"school_number,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	Phone         string `json:"phone,omitempty"`
	PhoneVerified bool   `json:"phone_verified,omitempty"`
}

func Callback(conf Config, code string) (*ResponseTokens, error) {
	if len(conf.OidcProviderURL) == 0 {
		conf.OidcProviderURL = "https://api.bingyan.net/sso/oidc"
	}
	provider, err := oidc.NewProvider(context.Background(), conf.OidcProviderURL)
	if err != nil {
		return nil, err
	}

	oauth2Config := &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  conf.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token found")
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: conf.ClientID})
	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return nil, err
	}

	var claims IDTokenClaims

	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return &ResponseTokens{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		IDToken:      rawIDToken,
		Scope:        token.Extra("scope").(string),
		ExpiresAt:    claims.ExpiresAt,
		RefreshToken: token.RefreshToken,
		Claims:       claims,
	}, nil
}
