package server

import (
	"encoding/json"
	"errors"
	"libsysfo-server/database"
	"libsysfo-server/utility/cred"
	"net/http"
	"time"
)

func profileBorrow(w http.ResponseWriter, r *http.Request) {
	data, invalid := checkIfAllowed(3, w, r)
	if invalid {
		return
	}

	borrowData := searchBorrow(data.ID)

	response{
		Data: responseBody{
			Borrow: borrowData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func searchBorrow(id int) []profileCollectionBorrowResponse {

	data := []database.LibraryCollectionBorrow{}
	database.DB.Where("user_id = ?", id).
		Preload("Collection.Library").
		Preload("Collection.Book").
		Preload("User").
		Order("created_at desc").
		Find(&data)

	return appendBorrowData(data)
}

func borrowNewBook(w http.ResponseWriter, r *http.Request) {
	var e newBorrowRequest
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

	collectionData := database.LibraryCollection{}
	err = database.DB.
		Where("id = ? ", e.Id).
		Find(&collectionData).
		Error
	if err != nil {
		intServerError(w, err)
	}

	if collectionData.Availability == 3 {
		badRequest(w, "Book is not available")
		return
	} else if collectionData.Availability == 2 {
		badRequest(w, "Book only available for read on library")
		return
	}

	borrowData := database.LibraryCollectionBorrow{
		UserID:       user.ID,
		CollectionID: e.Id,
		CreatedAt:    time.Now(),
	}

	err = database.DB.Save(&borrowData).Error

	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Borrow requested",
	}.responseFormatter(w)
}
