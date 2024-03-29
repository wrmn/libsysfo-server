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

func libraryBorrow(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	borrowsData := []profileCollectionBorrowResponse{}
	collectionsData := []database.LibraryCollection{}
	err := database.DB.Where("library_id = ?", libraryData.ID).
		Preload("Borrow", func(db *gorm.DB) *gorm.DB {
			return database.DB.
				Preload("Collection.Library").
				Preload("Collection.Book").
				Preload("User.ProfileData")
		}).
		Find(&collectionsData).Error

	if err != nil {
		intServerError(w, err)
		return
	}

	for _, k := range collectionsData {
		borrowsData = append(borrowsData, appendBorrowData(k.Borrow)...)
	}

	response{
		Data: responseBody{
			Borrow: &borrowsData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}

func findBorrow(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	q := r.URL.Query()
	if !(q.Has("cid") && q.Has("uid")) {
		badRequest(w, "Incomplete Parameters")
		return
	}

	cid, err := strconv.Atoi(q.Get("cid"))
	if err != nil {
		badRequest(w, "invalid parameter")
		return
	}

	uid, err := strconv.Atoi(q.Get("uid"))
	if err != nil {
		badRequest(w, "invalid parameter")
		return
	}

	borrowsData := []profileCollectionBorrowResponse{}
	collectionData := database.LibraryCollection{}
	db := database.DB.
		Where("library_id = ? AND id = ?", libraryData.ID, cid).
		Preload("Borrow", func(db *gorm.DB) *gorm.DB {
			return database.DB.
				Preload("Collection.Library").
				Preload("Collection.Book").
				Preload("User.ProfileData").
				Where("user_id = ? ", uid)
		}).
		Find(&collectionData)

	if invalid := databaseException(w, db); invalid {
		return
	}

	for _, l := range collectionData.Borrow {
		borrowsData = append(borrowsData, appendBorrowData([]database.LibraryCollectionBorrow{l})...)
	}

	userData, invalid := findUserById(uid, w)
	if invalid {
		return
	}

	resultCollection, invalid := findCollectionById(cid, w)
	if invalid {
		return
	}

	response{
		Data: responseBody{
			User: generateProfileResponse(userData),
			Book: setBookResponse(resultCollection.Book),
			Collection: libraryCollectionResponse{
				Id:           resultCollection.ID,
				SerialNumber: resultCollection.SerialNumber,
				Availability: resultCollection.Availability,
				Status:       resultCollection.Status,
			},
			Borrow: &borrowsData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)

}

func libraryUserBorrow(w http.ResponseWriter, r *http.Request) {
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
		Preload("Borrow.Collection.Book").
		Preload("Borrow.Collection.Library").
		Preload("Borrow.User.ProfileData").
		Preload("ProfileData").
		Find(&userResult).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	borrowData := appendBorrowData(userResult.Borrow)
	respBorrow := []profileCollectionBorrowResponse{}

	for _, i := range borrowData {
		if i.LibraryId == libraryData.ID {
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

func actionBorrow(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	var e borrowRequest
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

	if e.State == "next" || e.State == "cancel" {
		var borrow database.LibraryCollectionBorrow
		if e.BorrowId == nil {
			badRequest(w, "borrow Id required")
			return
		}

		db := database.DB.Where("id = ?", *e.BorrowId).
			Preload("Collection").
			Find(&borrow)
		if invalid := databaseException(w, db); invalid {
			return
		}

		if borrow.Collection.LibraryID != libraryData.ID {
			badRequest(w, "borrow Data not found")
			return
		}

		borrowNextAction(w, e)

	} else if e.State == "new" || e.State == "newTake" {
		borrowNewAction(w, e, libraryData.ID)
	} else {
		badRequest(w, "invalid state. use new, newTake, next, or cancel as state.")
	}
}

func borrowNewAction(w http.ResponseWriter, e borrowRequest, libraryId int) {
	if e.CollectionId == nil || e.UserId == nil {
		badRequest(w, "incomplete request, select user and collection corectly")
		return
	}

	collectionData := database.LibraryCollection{}
	db := database.DB.Where("id = ?", e.CollectionId).Find(&collectionData)
	if invalid := databaseException(w, db); invalid {
		return
	}
	if collectionData.LibraryID != libraryId {
		badRequest(w, "borrow Data not found")
		return
	}
	newBorrow := database.LibraryCollectionBorrow{
		CreatedAt:    time.Now(),
		CollectionID: *e.CollectionId,
		UserID:       *e.UserId,
	}
	newBorrow.AcceptedAt = &newBorrow.CreatedAt

	if e.State == "newTake" {
		newBorrow.TakedAt = &newBorrow.CreatedAt
	}

	if collectionData.Availability != 1 {
		badRequest(w, "Book is not available to borrow")
		return
	}

	err := database.DB.Create(&newBorrow).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	err = cancelOtherBorrow(newBorrow)
	if err != nil {
		intServerError(w, err)
		return
	}

	collectionData.Availability = 3
	err = database.DB.Save(&collectionData).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	err = database.DB.Create(&database.Notification{
		UserID:  newBorrow.UserID,
		Message: "New borrow has been added by admin",
		Read:    false,
	}).Error

	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "New Borrow created",
	}.responseFormatter(w)
}

func borrowNextAction(w http.ResponseWriter, e borrowRequest) {
	borrowData := database.LibraryCollectionBorrow{}
	database.DB.Where("id = ?", e.BorrowId).
		Preload("Collection").
		Find(&borrowData)

	now := time.Now()
	var msg string
	var notifMsg string
	if e.State == "next" {
		if borrowData.AcceptedAt == nil {
			err := cancelOtherBorrow(borrowData)
			if err != nil {
				intServerError(w, err)
				return
			}
			borrowData.AcceptedAt = &now
			borrowData.Collection.Availability = 3
			notifMsg = "borrow request has been accepted"
		} else if borrowData.TakedAt == nil {
			borrowData.TakedAt = &now
			borrowData.Collection.Availability = 3
			notifMsg = "borrow request has been taked"
		} else if borrowData.ReturnedAt == nil {
			borrowData.ReturnedAt = &now
			borrowData.Collection.Availability = 1
			notifMsg = "borrow request has been returned"
		} else {
			badRequest(w, "Borrow status already finished")
			return
		}
		msg = "Borrow status updated"
	} else if e.State == "cancel" && borrowData.CanceledAt == nil {
		borrowData.CanceledAt = &now
		borrowData.Collection.Availability = 1
		msg = "Borrow is rejected"
		notifMsg = "borrow request has been rejected"
	} else {
		badRequest(w, "Borrow status already canceled")
		return
	}

	err := database.DB.Save(&borrowData).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	err = database.DB.Save(&borrowData.Collection).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	err = database.DB.Create(&database.Notification{
		UserID:  borrowData.UserID,
		Message: notifMsg,
		Read:    false,
	}).Error

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
