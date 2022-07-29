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

	_, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	userData, err := findUserById(id)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	response{
		Data: responseBody{
			Profile: generateProfileResponse(userData),
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func libraryUserFind(w http.ResponseWriter, r *http.Request) {

	_, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	result := []profileResponse{}
	userResult := []database.ProfileAccount{}
	database.DB.Table("profile_accounts").
		Scopes(userFindFilter(r)).
		Preload("ProfileData").
		Find(&userResult)

	for i, k := range userResult {
		result = append(result, profileResponse{
			Id:       &userResult[i].ID,
			Username: k.Username,
			Email:    k.Email,
			Verified: k.ProfileData.VerifiedAt,
			Name:     k.ProfileData.Name,
		})
	}

	response{
		Data: responseBody{
			User: result,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}
