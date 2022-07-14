package cred

import (
	"libsysfo-server/database"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(data database.ProfileAccount) (t string, e error) {
	key := getKey(data.AccountType)
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
	t, e = token.SignedString([]byte(key))
	return
}

func VerifyToken(tokenString string) (t *jwt.Token, e error) {
	t, e = jwt.ParseWithClaims(tokenString, &TokenModel{}, func(token *jwt.Token) (interface{}, error) {
		cred := token.Claims.(*TokenModel)
		return []byte(getKey(cred.AccountType)), nil
	})
	return
}

func getKey(accType int) (key string) {
	if accType == 2 {
		key = os.Getenv("TOKEN_KEY_LIBRARY")
	} else if accType == 3 {
		key = os.Getenv("TOKEN_KEY")
	}
	return
}
