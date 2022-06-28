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
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	token, err := cred.VerifyGoogleToken(e.Credential)

	if err != nil {
		unauthorizedRequest(w, err)
	} else {
		user := database.ProfileAccount{}
		result := database.DB.Where("email = ?", token.Claims["email"]).Find(&user)
		if result.RowsAffected == 0 {
			err := googleRegisterHandler(token.Claims)
			if err != nil {
				unauthorizedRequest(w, err)
			}
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
		unauthorizedRequest(w, err)
	}
	response{
		Data:        tokenResult,
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}

func emailValidate(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token, present := query["token"]
	if !present || len(token) == 0 {
		badRequest(w, "token not present")
		return
	} else if len(token) > 1 {
		badRequest(w, "multiple token detected")
		return
	}
	tokenData, err := cred.VerifyToken(token[0])
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	profile := database.ProfileAccount{}
	err = database.DB.Preload("ProfileData").
		Where("email = ?", tokenData.Claims.(*cred.TokenModel).Email).
		Find(&profile).Error
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	if profile.ProfileData.VerifiedAt == nil {
		now := time.Now()
		profile.ProfileData.VerifiedAt = &now
		err = database.DB.Save(&profile.ProfileData).Error
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		http.Redirect(w, r, "http://localhost:3000", http.StatusSeeOther)
	} else {
		badRequest(w, "user already verified")
	}
}

func googleRegisterHandler(data map[string]interface{}) (err error) {
	lastAcc := database.ProfileAccount{}
	database.DB.Last(&lastAcc)
	password := gofakeit.Gamertag()
	user := database.ProfileAccount{
		Id:          lastAcc.Id + 1,
		Email:       data["email"].(string),
		AccountType: 3,
		Password:    fmt.Sprintf("%x", md5.Sum([]byte(password))),
	}

	user.ProfileData = database.ProfileData{
		UserId:     lastAcc.Id + 1,
		Name:       data["name"].(string),
		IsWhatsapp: false,
		Images:     data["picture"].(string),
	}

	tokenResult, err := cred.CreateToken(user)
	if err != nil {
		return
	}

	database.DB.Create(&user)

	link := fmt.Sprintf("http://localhost:5000/profile/validate?token=%s", tokenResult)
	fmt.Println(link)

	content := fmt.Sprintf("<html><head></head><body><p>Hello,</p><p>Password Sementara anda adalah %s</p>segera menuju <a href='%s'>link</a> ini untuk verifikasi akun anda</body>	</html>",
		password,
		link,
	)

	emailBody := email.Content{
		Subject:     "Email Verification",
		HtmlContent: content,
	}

	receiver := email.ToData{
		Name:  data["name"].(string),
		Email: data["email"].(string),
	}

	err = emailBody.SendEmail(receiver)
	return
}
