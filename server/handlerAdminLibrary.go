package server

import (
	"errors"
	"libsysfo-server/database"
	"libsysfo-server/utility/cred"
	"net/http"
)

func adminLogin(w http.ResponseWriter, r *http.Request) {
	user, err := getLoginData(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	} else if user.AccountType != 2 {
		err := errors.New("user not allowed")
		unauthorizedRequest(w, err)
		return
	}
	loginHandler(w, user)
}

func adminInformation(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	data := database.ProfileAccount{}
	cred := tokenData.Claims.(*cred.TokenModel)
	database.DB.Where("email = ?", cred.Email).Or("username = ?", cred.Username).
		Preload("Library").First(&data)

	response{
		Data: responseBody{
			Profile: adminInformationResponse{
				Username: *data.Username,
				Email:    data.Email,
				Library:  data.Library.Name,
				Image:    data.Library.ImagesMain,
				Address:  data.Library.Address,
			},
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}
