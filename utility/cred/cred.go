package cred

import (
	"libsysfo-server/database"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(data database.ProfileAccount) (t string, e error) {
	claims := TokenModel{
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

func VerifyToken(tokenString string) (t *jwt.Token, e error) {
	t, e = jwt.ParseWithClaims(tokenString, &TokenModel{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_KEY")), nil
	})
	return
}
