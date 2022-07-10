package server

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"libsysfo-server/database"
	"libsysfo-server/utility/cred"
	"libsysfo-server/utility/email"
	"libsysfo-server/utility/imgkit"
	"net/http"
	"strings"
)

func updatePassword(w http.ResponseWriter, r *http.Request) {
	var e profilePwdUpdateRequest
	var unmarshalErr *json.UnmarshalTypeError
	tokenData, err := authVerification(r)

	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	cred := tokenData.Claims.(*cred.TokenModel)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	if e.Password != e.RetypePassword {
		badRequest(w, "incorrect retype password")
		return
	}

	e.Password = fmt.Sprintf("%x", md5.Sum([]byte(e.Password)))
	user, err := findUser(cred, e.OldPassword)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	user.Password = e.Password
	database.DB.Save(&user)

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Password changed",
	}.responseFormatter(w)
}

func updateEmail(w http.ResponseWriter, r *http.Request) {
	var e profileEmailUpdateRequest
	var unmarshalErr *json.UnmarshalTypeError
	tokenData, err := authVerification(r)

	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	cred := tokenData.Claims.(*cred.TokenModel)
	pwd, err := pwdLocator(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	user, err := findUser(cred, pwd)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	user.Email = e.Email
	err = database.DB.Save(&user).Error
	if err != nil {
		badRequest(w, "email has been used")
		return
	}

	loginHandler(w, user)
}

func updateUsername(w http.ResponseWriter, r *http.Request) {
	var e profileUsernameUpdateRequest
	var unmarshalErr *json.UnmarshalTypeError
	tokenData, err := authVerification(r)

	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	cred := tokenData.Claims.(*cred.TokenModel)
	pwd, err := pwdLocator(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	user, err := findUser(cred, pwd)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	user.Username = e.Username
	err = database.DB.Save(&user).Error
	if err != nil {
		badRequest(w, "username has been used")
		return
	}

	loginHandler(w, user)
}

func updatePicture(w http.ResponseWriter, r *http.Request) {
	var e profilePictureUpdateRequest
	var unmarshalErr *json.UnmarshalTypeError
	tokenData, err := authVerification(r)

	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	cred := tokenData.Claims.(*cred.TokenModel)
	pwd, err := pwdLocator(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	user, err := findUser(cred, pwd)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	img := imgkit.ImgInformation{
		File:     e.Picture,
		FileName: "profile",
		Folder:   fmt.Sprintf("/user/%d/", user.ID),
	}

	upr, err := img.UploadImage()
	if err != nil {
		intServerError(w, err)
		return
	}

	userData, err := findUserData(user.ID)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	userData.Images = upr.URL

	err = database.DB.Save(&userData).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Picture changed",
	}.responseFormatter(w)
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	var e profileUpdateRequest
	var unmarshalErr *json.UnmarshalTypeError
	tokenData, err := authVerification(r)

	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	cred := tokenData.Claims.(*cred.TokenModel)
	pwd, err := pwdLocator(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	user, err := findUser(cred, pwd)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	userData, err := findUserData(user.ID)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	userData.Name = e.Name
	userData.Gender = e.Gender
	userData.PlaceOfBirth = e.PlaceOfBirth
	userData.DateOfBirth = e.DateOfBirth
	userData.Address1 = e.Address
	userData.Profession = e.Profession
	userData.Institution = e.Institution
	userData.PhoneCode = e.PhoneCode
	userData.PhoneNo = e.PhoneNo
	userData.IsWhatsapp = e.IsWhatsapp

	err = database.DB.Save(&userData).Error
	if err != nil {
		intServerError(w, err)
		return
	}
	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Profile changed",
	}.responseFormatter(w)
}

func resendEmail(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	tokenResult := strings.Split(r.Header.Values("Authorization")[0], " ")
	//NOTE: change to deployed url server
	link := fmt.Sprintf("http://localhost:5000/profile/validate?token=%s", tokenResult[1])
	fmt.Println(link)
	content := fmt.Sprintf("<html><head></head><body><p>Hello,</p>Segera menuju <a href='%s'>link</a> ini untuk verifikasi akun anda</body>	</html>",
		link,
	)

	emailBody := email.Content{
		Subject:     "Email Verification",
		HtmlContent: content,
	}

	receiver := email.ToData{
		Name:  "Libsysfo user",
		Email: tokenData.Claims.(*cred.TokenModel).Email,
	}

	err = emailBody.SendEmail(receiver)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Email send",
	}.responseFormatter(w)
}
