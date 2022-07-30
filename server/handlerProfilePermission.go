package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"libsysfo-server/database"
	"libsysfo-server/utility/cred"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func profileAccessPermission(w http.ResponseWriter, r *http.Request) {
	data, invalid := checkIfAllowed(3, w, r)
	if invalid {
		return
	}

	permissionData, err := searchPermission(data.ID)
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Data: responseBody{
			Permission: permissionData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func profileReadPaper(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	data, invalid := checkIfAllowed(3, w, r)
	if invalid {
		return
	}

	permissionData := database.LibraryPaperPermission{}
	db := database.DB.
		Where("id = ? AND user_id = ?", id, data.ID).
		Preload("Paper").Find(&permissionData)

	if invalid := databaseException(w, db); invalid {
		return
	}

	resp, err := http.Get(permissionData.Paper.PaperUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		intServerError(w, errors.New("drive server timeout"))
		return
	}

	access := database.LibraryPaperAccess{
		PermissionID: permissionData.ID,
		CreatedAt:    time.Now(),
	}

	database.DB.Save(&access)

	w.Write(body)
}

func searchPermission(id int) (respBody []profilePermissionResponse, err error) {
	data := []database.LibraryPaperPermission{}
	err = database.DB.Where("user_id = ?", id).
		Preload("Paper.Library").Find(&data).Error
	if err != nil {
		return
	}

	respBody = appendPermissionData(data)

	return
}

func profileNewPermission(w http.ResponseWriter, r *http.Request) {
	var e newPermissionRequest
	var unmarshalErr *json.UnmarshalTypeError
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	cred := tokenData.Claims.(*cred.TokenModel)
	user, err := getUser(cred)
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

	stats, err := userStats(user.ID)
	if err != nil || !stats {
		badRequest(w, "Please verify your email address and fill all profile information")
		return
	}

	statsPaper := database.DB.
		Where("id = ? AND access = true", e.Id).
		Find(&database.LibraryPaper{}).RowsAffected

	if statsPaper == 0 {
		badRequest(w, "Dokumen tidak mengizinkan akses")
		return
	}

	statsPermission := database.DB.
		Where("paper_id = ? AND user_id = ? AND accepted_at IS NOT NULL", e.Id, user.ID).
		Find(&database.LibraryPaperPermission{}).RowsAffected

	if statsPermission != 0 {
		response{
			Status:      http.StatusOK,
			Reason:      "Ok",
			Description: "Paper access has been requested before",
		}.responseFormatter(w)
		return
	}

	permission := database.LibraryPaperPermission{
		PaperID: e.Id,
		UserID:  user.ID,
		Purpose: e.Purpose,
	}

	err = database.DB.Save(&permission).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Paper access requested",
	}.responseFormatter(w)
}
