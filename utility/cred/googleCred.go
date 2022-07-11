package cred

import (
	"context"
	"libsysfo-server/database"
	"os"

	"google.golang.org/api/idtoken"
)

// verify google login token
func VerifyGoogleToken(idToken string) (payload *idtoken.Payload, err error) {
	payload, err = idtoken.Validate(context.Background(), idToken, os.Getenv("GOOGLE_CLIENT_ID"))
	database.DB.Save(&database.ThirdPartyJobs{
		Job:          "Verify Google Account",
		Destination:  "Google",
		ResponseBody: payload.Subject,
		Status:       200,
	})
	return
}
