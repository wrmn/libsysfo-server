package server

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"libsysfo-server/database"
	"libsysfo-server/utility/cred"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func findUser(cred *cred.TokenModel, pwd string) (user database.ProfileAccount, err error) {
	password := fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	result := database.DB.
		Where("email = ? AND password = ?", cred.Email, password).
		Or("username = ? AND password = ?", cred.Username, password).
		Find(&user)

	if result.RowsAffected == 0 {
		err = errors.New("invalid password")
		return
	}
	return
}

func getUser(cred *cred.TokenModel) (user database.ProfileAccount, err error) {
	err = database.DB.
		Where("email = ? ", cred.Email).
		Or("username = ? ", cred.Username).
		Find(&user).Error
	if user.AccountType != 3 {
		err = errors.New("user not allowed")
		return
	}
	return
}

func checkToken(r *http.Request, w http.ResponseWriter) (database.ProfileAccount, bool) {
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return database.ProfileAccount{}, true
	}
	data, err := userData(tokenData)
	if err != nil {
		unauthorizedRequest(w, err)
		return database.ProfileAccount{}, true
	}
	return data, false
}

func userData(tokenData *jwt.Token) (data database.ProfileAccount, err error) {
	cred := tokenData.Claims.(*cred.TokenModel)
	db := database.DB.Where("email = ?", cred.Email).Or("username = ?", cred.Username).
		Preload("Library").Preload("ProfileData").First(&data)
	if db.RowsAffected != 1 {
		return data, errors.New("user not found")
	}
	return data, db.Error
}

func authVerification(r *http.Request) (tokenData *jwt.Token, err error) {
	tokenHeader := r.Header.Values("Authorization")
	if len(tokenHeader) == 0 {
		err = errors.New("authorization required")
		return
	}
	token := strings.Split(tokenHeader[0], " ")
	if token[0] != "Bearer" {
		err = errors.New("need bearer authorization")
		return
	}
	tokenData, err = cred.VerifyToken(token[1])
	return
}

func pwdLocator(r *http.Request) (pwd string, err error) {
	pwdHead := r.Header.Values("Account-auth")
	if len(pwdHead) == 0 {
		err = errors.New("authorization required")
		return
	}
	pwd = pwdHead[0]
	return
}

func findUserById(id int) (data database.ProfileAccount, err error) {
	result := database.DB.
		Where("id = ?", id).
		Preload("ProfileData").
		Find(&data)

	if result.RowsAffected == 0 {
		err = errors.New("user not found")
		return
	}

	err = result.Error
	return
}

func findUserData(id int) (user database.ProfileData, err error) {
	result := database.DB.
		Where("user_id = ?", id).
		Find(&user)

	if result.RowsAffected == 0 {
		err = errors.New("user not found")
		return
	}
	return
}

func userStats(id int) (stats bool, err error) {
	res, err := findUserData(id)
	if res.Address1 == nil ||
		res.DateOfBirth == nil ||
		res.Gender == nil ||
		res.Institution == nil ||
		res.PhoneCode == nil ||
		res.PlaceOfBirth == nil ||
		res.Profession == nil ||
		res.VerifiedAt == nil {
		stats = false
		return
	}
	return true, err
}

func getLoginData(r *http.Request) (user database.ProfileAccount, err error) {
	var e cred.FormAuth
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()
	err = decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			err = errors.New("Wrong Type provided for field " + unmarshalErr.Field)
		}
		return
	}

	e.Password = fmt.Sprintf("%x", md5.Sum([]byte(e.Password)))

	user = database.ProfileAccount{}
	result := database.DB.Where("email = ? AND password = ?", e.Indicator, e.Password).Or("username = ? AND password = ?", e.Indicator, e.Password).Find(&user).RowsAffected
	if result == 0 {
		err = errors.New("invalid username or password ")
	}

	return
}

func userFindFilter(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		db = db.Where("account_type = ?", 3)
		if q.Has("keyword") {
			keyword := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("keyword")))
			db = db.Where("LOWER(username) like ? AND account_type = ?", keyword, 3).
				Or("LOWER(email) like ? AND account_type = ?", keyword, 3)
		}
		return db
	}
}

func generateProfileResponse(data database.ProfileAccount) profileResponse {
	return profileResponse{
		Username:     data.Username,
		Email:        data.Email,
		Verified:     data.ProfileData.VerifiedAt,
		Name:         data.ProfileData.Name,
		Gender:       data.ProfileData.Gender,
		PlaceOfBirth: data.ProfileData.PlaceOfBirth,
		DateOfBirth:  data.ProfileData.DateOfBirth,
		Address:      data.ProfileData.Address1,
		Institution:  data.ProfileData.Institution,
		Profession:   data.ProfileData.Profession,
		PhoneCode:    data.ProfileData.PhoneCode,
		PhoneNo:      data.ProfileData.PhoneNo,
		IsWhatsapp:   &data.ProfileData.IsWhatsapp,
		Images:       data.ProfileData.Images,
	}
}
