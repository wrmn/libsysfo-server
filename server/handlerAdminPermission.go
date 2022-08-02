package server

import (
	"encoding/json"
	"errors"
	"libsysfo-server/database"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func libraryPermission(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	data := []database.LibraryPaper{}
	err := database.DB.
		Where("library_id = ?", libraryData.ID).
		Preload("Permission.Paper.Library").
		Preload("Permission.User.ProfileData").
		Find(&data).
		Error

	if err != nil {
		intServerError(w, err)
		return
	}

	respBody := []profilePermissionResponse{}
	for _, d := range data {
		respBody = append(respBody, appendPermissionData(d.Permission)...)
	}

	response{
		Data: responseBody{
			Permission: respBody,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func findPermission(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	q := r.URL.Query()
	if !(q.Has("pid") && q.Has("uid")) {
		badRequest(w, "Incomplete Parameters")
		return
	}

	pid, err := strconv.Atoi(q.Get("pid"))
	if err != nil {
		badRequest(w, "invalid parameter")
		return
	}

	uid, err := strconv.Atoi(q.Get("uid"))
	if err != nil {
		badRequest(w, "invalid parameter")
		return
	}

	paperData := database.LibraryPaper{}

	db := database.DB.
		Where("library_id = ? AND id=?", libraryData.ID, pid).
		Preload("Permission", func(db *gorm.DB) *gorm.DB {
			return database.DB.
				Preload("Access").
				Where("user_id = ?", uid)
		}).Find(&paperData)
	if invalid := databaseException(w, db); invalid {
		return
	}
	if len(paperData.Permission) < 1 {
		badRequest(w, "permission not found")
		return
	}
	permission := paperData.Permission[0]
	permissionData := formatPermissionData(permission)

	userData, invalid := findUserById(uid, w)
	if invalid {
		return
	}

	paperResponse := setPaperResponse(paperData)
	paperResponse.PaperUrl = &paperData.PaperUrl

	respBody := accessResponse{
		Total:     len(permission.Access),
		CreatedAt: appendAccessData(permission.Access),
	}

	response{
		Data: responseBody{
			User:       generateProfileResponse(userData),
			Paper:      paperResponse,
			Permission: &permissionData,
			Access:     &respBody,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}

func accessHistory(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "invalid id")
		return
	}

	permissionData := database.LibraryPaperPermission{}

	db := database.DB.Where("id = ?", id).
		Preload("Paper").
		Preload("Access").
		Preload("User.ProfileData").Find(&permissionData)

	if invalid := databaseException(w, db); invalid {
		return
	}

	if permissionData.Paper.LibraryID != libraryData.ID {
		badRequest(w, "user not allowed")
		return
	}

	respBody := accessResponse{
		Total:     len(permissionData.Access),
		CreatedAt: appendAccessData(permissionData.Access),
	}

	response{
		Data: responseBody{
			Access: respBody,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}

func actionPermission(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	now := time.Now()
	var msg string
	var e permissionRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	permissionData := database.LibraryPaperPermission{}

	db := database.DB.Where("id = ?", e.PermissionId).
		Preload("Paper").Find(&permissionData)

	if invalid := databaseException(w, db); invalid {
		return
	}

	if permissionData.Paper.LibraryID != libraryData.ID {
		badRequest(w, "user not allowed")
		return
	}

	if permissionData.AcceptedAt != nil || permissionData.CanceledAt != nil {
		badRequest(w, "permission request has been responded")
		return
	}

	if e.State == "accept" {
		permissionData.AcceptedAt = &now
		msg = "request accepted"
	} else if e.State == "cancel" {
		permissionData.CanceledAt = &now
		msg = "request rejected"
	} else {
		badRequest(w, "invalid state. use accept or cancel")
		return
	}
	err = database.DB.Save(&permissionData).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: msg,
	}.responseFormatter(w)
}

func libraryUserPermission(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	userResult := database.ProfileAccount{}
	err = database.DB.Where("id = ?", id).
		Preload("Permission.Paper.Library").
		Preload("Permission.User.ProfileData").
		Preload("ProfileData").
		Find(&userResult).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	permissionData := appendPermissionData(userResult.Permission)
	respPermission := []profilePermissionResponse{}

	for _, i := range permissionData {
		if i.LibraryId == libraryData.ID {
			respPermission = append(respPermission, i)
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
			Permission: &respPermission,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}
