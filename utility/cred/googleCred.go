package cred

import (
	"context"
	"os"

	"google.golang.org/api/idtoken"
)

// verify google login token
func VerifyGoogleToken(idToken string) (payload *idtoken.Payload, err error) {
	payload, err = idtoken.Validate(context.Background(), idToken, os.Getenv("GOOGLE_CLIENT_ID"))
	return
}
