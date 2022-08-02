package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"libsysfo-server/database"
	"libsysfo-server/utility/imgkit"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func libraryPaper(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	respBody := []paperResponse{}
	paperData, err := getLibraryPaper(libraryData.ID)
	if err != nil {
		intServerError(w, err)
	}

	response{
		Data: responseBody{
			Paper: append(respBody, paperData...),
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}

func librarySinglePaper(w http.ResponseWriter, r *http.Request) {

	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	paperId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "invalid id request")
	}

	result := database.LibraryPaper{}
	err = database.DB.
		Where("Id = ? AND library_id = ?", paperId, libraryData.ID).
		Preload("Permission.User.ProfileData").
		Find(&result).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	paperResponse := setPaperResponse(result)
	paperResponse.PaperUrl = &result.PaperUrl

	permissionResponse := appendPermissionData(result.Permission)

	response{
		Data: responseBody{
			Paper:      paperResponse,
			Permission: permissionResponse,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}

func libraryAddPaper(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}
	var e paperAddRequest
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

	paperData := database.LibraryPaper{
		LibraryID:   libraryData.ID,
		Title:       e.Title,
		Subject:     e.Subject,
		Abstract:    e.Abstract,
		Type:        e.Type,
		Description: e.Description,
		PaperUrl:    "https://ik.imagekit.io/libsysfo/test.pdf",
	}

	if e.Access != nil {
		paperData.Access = *e.Access
	} else {
		paperData.Access = false
	}

	err = database.DB.Create(&paperData).Error

	if err != nil {
		intServerError(w, err)
		return
	}

	file := imgkit.ImgInformation{
		File:     e.PaperFile,
		FileName: strconv.Itoa(paperData.ID),
		Folder:   fmt.Sprintf("/paper/%d/", paperData.ID),
	}
	upr, err := file.UploadImage()
	if err != nil {
		intServerError(
			w,
			fmt.Errorf("data added to database, but file fail to upload with error: %s, please edit file to upload file again ",
				err.Error()))
		return
	}

	paperData.PaperUrl = upr.URL

	err = database.DB.Save(&paperData).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Paper Saved",
	}.responseFormatter(w)
}

func libraryUpdatePaper(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	paperId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "invalid id request")
	}

	result := database.LibraryPaper{}
	err = database.DB.Where("Id = ? AND library_id = ?", paperId, libraryData.ID).Find(&result).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	var e paperAddRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	result.Title = e.Title
	result.Subject = e.Subject
	result.Access = *e.Access
	result.Abstract = e.Abstract
	result.Type = e.Type
	result.Description = e.Description

	err = database.DB.Save(&result).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Paper updated",
	}.responseFormatter(w)
}

func libraryUpdatePaperFile(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	paperId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "invalid id request")
	}

	result := database.LibraryPaper{}
	err = database.DB.Where("Id = ? AND library_id = ?", paperId, libraryData.ID).Find(&result).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	var e fileUpdateRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	file := imgkit.ImgInformation{
		File:     e.File,
		FileName: strconv.Itoa(result.ID),
		Folder:   fmt.Sprintf("/paper/%d/", result.ID),
	}
	upr, err := file.UploadImage()
	if err != nil {
		intServerError(
			w,
			fmt.Errorf("data added to database, but file fail to upload with error: %s, please edit file to upload file again ",
				err.Error()))
		return
	}

	result.PaperUrl = upr.URL

	err = database.DB.Save(&result).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Paper Saved",
	}.responseFormatter(w)
}
