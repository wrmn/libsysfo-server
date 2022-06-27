package cred

import (
	"libsysfo-server/database"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(data database.ProfileAccount) (t string, e error) {
	claims := tokenModel{
		Username:    data.Username,
		Email:       data.Email,
		AccountType: data.AccountType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, e = token.SignedString([]byte(os.Getenv("TOKEN_KEY")))
	return
}
