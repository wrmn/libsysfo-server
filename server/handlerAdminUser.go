package server

import (
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
