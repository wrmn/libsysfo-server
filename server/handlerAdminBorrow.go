package server

import (
	"errors"
	"libsysfo-server/database"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func libraryBorrow(w http.ResponseWriter, r *http.Request) {

	data, invalid := checkToken(r, w)
	if invalid {
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
	err = database.DB.Where("library_id", libOwn.ID).
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

func getBorrow(w http.ResponseWriter, r *http.Request) {
	data, invalid := checkToken(r, w)
	if invalid {
		return
	}
	if data.AccountType != 2 {
		unauthorizedRequest(w, errors.New("user not allowed"))
		return
	}

	q := r.URL.Query()
	if !(q.Has("cid") && q.Has("uid")) {
		badRequest(w, "Incomplete Parameters")
		return
	}
	cid, _ := strconv.Atoi(q.Get("cid"))
	uid, _ := strconv.Atoi(q.Get("uid"))

	respBorrow := []database.LibraryCollectionBorrow{}
	db := database.DB.Where("collection_id = ? AND user_id = ?", cid, uid).
		Preload("Collection.Library").
		Preload("Collection.Book").
		Preload("User.ProfileData").
		Find(&respBorrow)

	if db.RowsAffected < 1 {
		return
	}
	if db.Error != nil {
		intServerError(w, db.Error)
	}

}
