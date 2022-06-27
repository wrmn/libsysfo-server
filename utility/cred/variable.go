package cred

import "github.com/golang-jwt/jwt/v4"

type GoogleAuth struct {
	ClientId   string `json:"clientId"`
	Credential string `json:"credential"`
	SelectBy   string `json:"select_by"`
}

type GoogleAuthClaims struct {
	jwt.RegisteredClaims
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Azp           string `json:"azp"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Jti           string `json:"jti"`
}

type tokenModel struct {
	jwt.RegisteredClaims
	Username    string `json:"username"`
	Email       string `json:"email"`
	AccountType int    `json:"accountType"`
}
