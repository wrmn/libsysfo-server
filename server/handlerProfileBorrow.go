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
	data, invalid := checkToken(r, w)
	if invalid {
		return
	}

	if data.AccountType != 3 {
		unauthorizedRequest(w, errors.New("user not allowed"))
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

func searchBorrow(id int) []profileCollectionBorrow {

	data := []database.LibraryCollectionBorrow{}
	database.DB.Where("user_id = ?", id).
		Preload("Collection.Library").
		Preload("Collection.Book").
		Preload("User").
		Order("created_at desc").
		Find(&data)

	return appendData(data)
}

func appendData(data []database.LibraryCollectionBorrow) (respBody []profileCollectionBorrow) {
	for _, d := range data {
		respBody = append(respBody, profileCollectionBorrow{
			BorrowId:     d.ID,
			CreatedAt:    d.CreatedAt,
			AcceptedAt:   d.AcceptedAt,
			TakedAt:      d.TakedAt,
			ReturnedAt:   d.ReturnedAt,
			CanceledAt:   d.CanceledAt,
			Title:        d.Collection.Book.Title,
			SerialNumber: d.Collection.SerialNumber,
			CollectionId: d.Collection.ID,
			Slug:         d.Collection.Book.Slug,
			LibraryId:    d.Collection.LibraryID,
			Library:      d.Collection.Library.Name,
			UserId:       d.User.ID,
			UserName:     d.User.ProfileData.Name,
			Status:       setStatus(d),
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
	statsCollection := database.DB.
		Where("id = ? AND availability = ?", e.Id, 1).
		Find(&collectionData).
		RowsAffected

	if statsCollection == 0 {
		badRequest(w, "Book is not available")
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
