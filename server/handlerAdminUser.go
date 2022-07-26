package server

import (
	"libsysfo-server/database"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func libraryUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "invalid id")
		return
	}

	_, invalid := checkToken(r, w)
	if invalid {
		return
	}

	data, err := findUserById(id)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	response{
		Data: responseBody{
			Profile: generateProfileResponse(data),
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func libraryUserBorrow(w http.ResponseWriter, r *http.Request) {
	data, invalid := checkToken(r, w)
	if invalid {
		return
	}

	libOwn := database.LibraryData{}

	err := database.DB.Where("user_id = ?", data.ID).Find(&libOwn).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	userResult := database.ProfileAccount{}
	err = database.DB.Where("id = ?", id).
		Preload("Borrow.Collection.Book").
		Preload("Borrow.Collection.Library").
		Preload("Borrow.User.ProfileData").
		Preload("ProfileData").
		Find(&userResult).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	borrowData := appendData(userResult.Borrow)
	respBorrow := []profileCollectionBorrow{}

	for _, i := range borrowData {
		if i.LibraryId == libOwn.ID {
			respBorrow = append(respBorrow, i)
		}
	}

	response{
		Data: responseBody{
			User: profileResponse{
				Id:       &userResult.ID,
				Username: userResult.Username,
				Email:    userResult.Email,
				Verified: userResult.ProfileData.VerifiedAt,
				Name:     userResult.ProfileData.Name,
			},
			Borrow: &respBorrow,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}
