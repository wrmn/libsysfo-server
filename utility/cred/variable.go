package cred

import "github.com/golang-jwt/jwt/v4"

type GoogleAuth struct {
	ClientId   string `json:"clientId"`
	Credential string `json:"credential"`
	SelectBy   string `json:"select_by"`
}

type FormAuth struct {
	Indicator string `json:"indicator"`
	Password  string `json:"password"`
}

type RegisForm struct {
	Username       *string `json:"username"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	RetypePassword string  `json:"retypePassword"`
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

type TokenModel struct {
	jwt.RegisteredClaims
	Username    *string `json:"username"`
	Email       string  `json:"email"`
	AccountType int     `json:"accountType"`
}
