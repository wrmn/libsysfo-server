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
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	data := database.ProfileAccount{}
	cred := tokenData.Claims.(*cred.TokenModel)
	database.DB.Where("email = ?", cred.Email).Or("username = ?", cred.Username).
		Preload("ProfileData").First(&data)

	permissionData := searchPermission(data.ID)

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
	tokenData, err := authVerification(r)
	id := mux.Vars(r)["id"]
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	data := database.ProfileAccount{}
	cred := tokenData.Claims.(*cred.TokenModel)
	database.DB.Where("email = ?", cred.Email).Or("username = ?", cred.Username).
		Preload("ProfileData").First(&data)

	permissionData := database.LibraryPaperPermission{}
	row := database.DB.
		Where("id = ? AND user_id = ?", id, data.ID).
		Preload("Paper").Find(&permissionData).RowsAffected

	if row == 0 {
		unauthorizedRequest(w, errors.New("data not found"))
	}

	resp, err := http.Get(permissionData.Paper.PaperUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	access := database.LibraryPaperAccess{
		PermissionID: permissionData.ID,
		CreatedAt:    time.Now(),
	}

	database.DB.Save(&access)

	w.Write(body)
}

func searchPermission(id int) (respBody []profilePermissionResponse) {
	data := []database.LibraryPaperPermission{}
	database.DB.Where("user_id = ?", id).
		Preload("Paper.Library").Find(&data)

	for _, d := range data {
		respBody = append(respBody, profilePermissionResponse{
			CreatedAt:    d.CreatedAt,
			Id:           d.ID,
			PaperTitle:   d.Paper.Title,
			PaperSubject: d.Paper.Subject,
			PaperType:    d.Paper.Type,
			Library:      d.Paper.Library.Name,
			Purpose:      d.Purpose,
			Accepted:     d.Accepted,
		})
	}
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

	user, err := getUser(cred)
	if err != nil {
		unauthorizedRequest(w, err)
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
		Where("paper_id = ? AND user_id = ?", e.Id, user.ID).
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
