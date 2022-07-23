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
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	data := database.ProfileAccount{}
	cred := tokenData.Claims.(*cred.TokenModel)
	database.DB.Where("email = ?", cred.Email).Or("username = ?", cred.Username).
		Preload("ProfileData").First(&data)

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

func searchBorrow(id int) (respBody []profileCollectionBorrow) {

	data := []database.LibraryCollectionBorrow{}
	database.DB.Where("user_id = ?", id).
		Preload("Collection.Library").
		Preload("Collection.Book").
		Order("created_at desc").
		Find(&data)

	for _, d := range data {
		respBody = append(respBody, profileCollectionBorrow{
			CreatedAt:    d.CreatedAt,
			TakedAt:      d.TakedAt,
			ReturnedAt:   d.ReturnedAt,
			Title:        d.Collection.Book.Title,
			SerialNumber: d.Collection.SerialNumber,
			Slug:         d.Collection.Book.Slug,
			LibraryId:    d.Collection.LibraryID,
			Library:      d.Collection.Library.Name,
			Status:       d.Status,
		})
	}
	return
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
	collectionData := database.LibraryCollection{}

	statsCollection := database.DB.
		Where("id = ? AND availability = true", e.Id).
		Find(&collectionData).
		RowsAffected

	if statsCollection == 0 {
		badRequest(w, "Buku tidak tersedia")
		return
	}

	collectionData.Availability = 0

	borrowData := database.LibraryCollectionBorrow{
		UserID:       user.ID,
		CollectionID: e.Id,
		CreatedAt:    time.Now(),
		Status:       "requested",
	}
	err = database.DB.Save(&collectionData).Error

	if err != nil {
		intServerError(w, err)
		return
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
