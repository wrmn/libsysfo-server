package server

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"libsysfo-server/database"
	"libsysfo-server/utility/cred"
	"libsysfo-server/utility/email"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func loginGoogle(w http.ResponseWriter, r *http.Request) {

	var e cred.GoogleAuth
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()
	err := decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			response{
				Status:      http.StatusBadRequest,
				Reason:      "Bad Request.",
				Description: "Wrong Type provided for field " + unmarshalErr.Field,
			}.responseFormatter(w)
		} else {
			response{
				Status:      http.StatusBadRequest,
				Reason:      "Bad Request",
				Description: err.Error(),
			}.responseFormatter(w)
		}
		return
	}

	token, err := cred.VerifyGoogleToken(e.Credential)

	user := database.ProfileAccount{}

	if err != nil {
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(err)
		response{
			Status:      http.StatusUnauthorized,
			Reason:      "Unauthorized",
			Description: err.Error(),
		}.responseFormatter(w)
	} else {
		result := database.DB.Where("email = ?", token.Claims["email"]).Find(&user)

		if result.RowsAffected == 0 {
			googleRegisterHandler(token.Claims)
			database.DB.Where("email = ?", token.Claims["email"]).Find(&user)
		}
		now := time.Now()
		user.LastLogin = &now
		database.DB.Save(&user)
		loginHandler(w, user)
	}
}

func loginHandler(w http.ResponseWriter, data database.ProfileAccount) {
	tokenResult, err := cred.CreateToken(data)
	if err != nil {
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(tokenResult)
}

func googleRegisterHandler(data map[string]interface{}) {
	lastAcc := database.ProfileAccount{}
	database.DB.Last(&lastAcc)
	password := gofakeit.Gamertag()
	user := database.ProfileAccount{
		Id:          lastAcc.Id + 1,
		Email:       data["email"].(string),
		AccountType: 3,
		Password:    fmt.Sprintf("%x", md5.Sum([]byte(password))),
	}

	content := fmt.Sprintf("<html><head></head><body><p>Hello,</p><p>Password Sementara anda adalah %s</p>segera menuju <a href='%s'>link</a> ini untuk verifikasi akun anda</body>	</html>",
		password,
		"google.com",
	)

	emailBody := email.Content{
		Subject:     "Email Verification",
		HtmlContent: content,
	}

	receiver := email.ToData{
		Name:  data["name"].(string),
		Email: data["email"].(string),
	}

	emailBody.SendEmail(receiver)

	database.DB.Create(&user)
}
