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

	data, invalid := checkToken(r, w)
	if invalid {
		return
	}

	if data.AccountType != 2 {
		unauthorizedRequest(w, errors.New("user not allowed"))
		return
	}

	libOwn := database.LibraryData{}

	err := database.DB.Where("user_id = ?", data.ID).Find(&libOwn).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	borrowsData := []profileCollectionBorrow{}
	collectionsData := []database.LibraryCollection{}
	err = database.DB.Where("library_id = ?", libOwn.ID).
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
		borrowsData = append(borrowsData, appendData(k.Borrow)...)
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
	data, invalid := checkToken(r, w)
	if invalid {
		return
	}
	if data.AccountType != 2 {
		unauthorizedRequest(w, errors.New("user not allowed"))
		return
	}

	libOwn := database.LibraryData{}

	err := database.DB.Where("user_id = ?", data.ID).Find(&libOwn).Error
	if err != nil {
		intServerError(w, err)
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

	borrowsData := []profileCollectionBorrow{}
	collectionsData := []database.LibraryCollection{}
	db := database.DB.Where("library_id = ?", libOwn.ID).
		Preload("Borrow", func(db *gorm.DB) *gorm.DB {
			return database.DB.
				Preload("Collection.Library").
				Preload("Collection.Book").
				Preload("User.ProfileData")
		}).
		Find(&collectionsData)

	if db.RowsAffected < 1 {
		return
	}
	if db.Error != nil {
		intServerError(w, db.Error)
		return
	}

	for _, k := range collectionsData {
		for _, l := range k.Borrow {
			if l.UserID == uid && l.CollectionID == cid {
				borrowsData = append(borrowsData, appendData([]database.LibraryCollectionBorrow{l})...)
			}
		}
	}

	userData, err := findUserById(uid)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	resultCollection, err := findCollectionById(cid)
	if err != nil {
		intServerError(w, err)
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
	data, invalid := checkToken(r, w)
	if invalid {
		return
	}

	if data.AccountType != 2 {
		unauthorizedRequest(w, errors.New("user not allowed"))
		return
	}

	libOwn := database.LibraryData{}

	err := database.DB.Where("user_id = ?", data.ID).Find(&libOwn).Error
	if err != nil {
		intServerError(w, err)
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

	borrowData := appendData(userResult.Borrow)
	respBorrow := []profileCollectionBorrow{}

	for _, i := range borrowData {
		if i.LibraryId == libOwn.ID {
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
	data, invalid := checkToken(r, w)
	if invalid {
		return
	}

	if data.AccountType != 2 {
		unauthorizedRequest(w, errors.New("user not allowed"))
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
	var borrow database.LibraryCollectionBorrow
	if e.State == "next" || e.State == "cancel" {
		if e.BorrowId == nil {
			badRequest(w, "borrow Id required")
			return
		}

		db := database.DB.Where("id = ?", *e.BorrowId).
			Preload("Collection").
			Find(&borrow)

		if db.RowsAffected < 1 ||
			db.Error != nil ||
			borrow.Collection.LibraryID != data.Library.ID {
			badRequest(w, "borrow Data not found")
			return
		}

		borrowNextAction(w, e)

	} else if e.State == "new" || e.State == "newTake" {

		if e.CollectionId == nil || e.UserId == nil {
			badRequest(w, "incomplete request, select user and collection corectly")
			return
		}

		collectionData := database.LibraryCollection{}
		err := database.DB.Where("id = ?", e.CollectionId).Find(&collectionData).Error
		if err != nil {
			badRequest(w, err.Error())
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

		err = database.DB.Create(&newBorrow).Error
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

		response{
			Status:      http.StatusOK,
			Reason:      "Ok",
			Description: "New Borrow created",
		}.responseFormatter(w)

	} else {
		badRequest(w, "invalid state. use new, newTake, next, or cancel as state.")
		return
	}
}

func borrowNextAction(w http.ResponseWriter, e borrowRequest) {
	borrowData := database.LibraryCollectionBorrow{}
	database.DB.Where("id = ?", e.BorrowId).
		Preload("Collection").
		Find(&borrowData)

	now := time.Now()
	var msg string

	if e.State == "next" {
		if borrowData.AcceptedAt == nil {
			borrowData.AcceptedAt = &now
			borrowData.Collection.Availability = 3
		} else if borrowData.TakedAt == nil {
			borrowData.TakedAt = &now
			borrowData.Collection.Availability = 3
		} else if borrowData.ReturnedAt == nil {
			borrowData.ReturnedAt = &now
			borrowData.Collection.Availability = 1
		} else {
			badRequest(w, "Borrow status already finished")
			return
		}
		msg = "Borrow status updated"
	} else if e.State == "cancel" && borrowData.CanceledAt == nil {
		borrowData.CanceledAt = &now
		borrowData.Collection.Availability = 1
		msg = "Borrow is rejected"
	} else {
		badRequest(w, "Borrow status already canceled")
		return
	}

	err := database.DB.Save(&borrowData).Error
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
