package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"libsysfo-server/database"
	"libsysfo-server/utility/imgkit"
	"net/http"

	"github.com/gorilla/mux"
)

func libraryBooks(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	data, err := adminData(tokenData)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	libOwn := database.LibraryData{}

	err = database.DB.Where("user_id = ?", data.ID).Find(&libOwn).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	bookData, err := getLibraryBook(libOwn.ID)
	if err != nil {
		intServerError(w, err)
		return
	}
	response{
		Data: responseBody{
			Book: &bookData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}

func librarySingleBook(w http.ResponseWriter, r *http.Request) {
	collectionId := mux.Vars(r)["id"]
	resultCollection := database.LibraryCollection{}
	resultBook := database.Book{}
	err := database.DB.Where("id = ?", collectionId).
		Find((&resultCollection)).Error
	if err != nil {
		intServerError(w, err)
	}

	err = database.DB.Preload("BookDetail").
		Where("id = ?", resultCollection.BookID).
		Find(&resultBook).Error
	if err != nil {
		intServerError(w, err)
	}
	response{
		Data: responseBody{
			Book: setBookResponse(resultBook),
			Collection: libraryCollectionResponse{
				Id:           resultCollection.ID,
				SerialNumber: resultCollection.SerialNumber,
				Availability: resultCollection.Availability,
				Status:       resultCollection.Status,
			},
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func libraryAddCollection(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	data, err := adminData(tokenData)
	if err != nil {
		unauthorizedRequest(w, err)
		fmt.Println("here")
		return
	}

	libOwn := database.LibraryData{}

	err = database.DB.Where("user_id = ?", data.ID).Find(&libOwn).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	var e collectionAddRequests
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
	result := database.Book{}
	if e.Book != nil {
		slug := slugGenerator(e.Book.Title)
		img := imgkit.ImgInformation{
			File:     e.Book.Image,
			FileName: "book_cover",
			Folder:   fmt.Sprintf("/book/%s/", slug),
		}

		upr, err := img.UploadImage()
		if err != nil {
			intServerError(w, err)
			return
		}

		result = database.Book{
			Image:  upr.URL,
			Title:  e.Book.Title,
			Author: e.Book.Author,
			Source: "local",
			Slug:   slug,
		}

		err = database.DB.Create(&result).Error
		if err != nil {
			intServerError(w, err)
			return
		}

		err = database.DB.Create(&database.BookDetail{
			ID:          result.ID,
			ReleaseDate: e.Book.ReleaseDate,
			Description: e.Book.Description,
			Language:    e.Book.Language,
			Country:     e.Book.Country,
			Publisher:   e.Book.Publisher,
			PageCount:   e.Book.PageCount,
			Category:    e.Book.Category,
		}).Error

		if err != nil {
			intServerError(w, err)
			return
		}
	} else {
		query := database.DB.Preload("BookDetail").
			Where("slug = ?", *e.BookSlug).Find(&result)
		exist, err := checkExist(query)
		if err != nil {
			intServerError(w, err)
			return
		}
		if exist < 1 {
			path := fmt.Sprintf("/api/books/%s/detail", *e.BookSlug)
			data, err := serverEndpoint(path)
			if err != nil {
				intServerError(w, err)
				return
			}
			book := data.Book
			details := book.Detail

			result = database.Book{
				Image:  *book.Image,
				Title:  *book.Title,
				Author: *book.Author,
				Source: "gramedia",
				Slug:   *book.Slug,
			}

			err = database.DB.Create(&result).Error
			if err != nil {
				intServerError(w, err)
				return
			}

			err = database.DB.Create(&database.BookDetail{
				ID:          result.ID,
				ReleaseDate: *details.ReleaseDate,
				Description: *details.Description,
				Language:    *details.Language,
				Country:     *details.Country,
				Publisher:   *details.Publisher,
				PageCount:   int(*details.PageCount),
				Category:    *details.Category,
			}).Error

			if err != nil {
				intServerError(w, err)
				return
			}
		}
	}

	for _, k := range e.Collection {
		err = database.DB.Create(&database.LibraryCollection{
			SerialNumber: k.SerialNumber,
			LibraryID:    libOwn.ID,
			BookID:       result.ID,
			Availability: k.Availability,
			Status:       1,
		}).Error
		if err != nil {
			intServerError(w, err)
			return
		}
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Collection saved",
	}.responseFormatter(w)
}

func libraryUpdateCollection(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	collectionId := mux.Vars(r)["id"]

	auth, err := adminData(tokenData)
	if err != nil {
		unauthorizedRequest(w, err)
		fmt.Println("here")
		return
	}

	var e collectionUpdateRequest
	var unmarshalErr *json.UnmarshalTypeError

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

	resultCollection := database.LibraryCollection{}

	err = database.DB.Where("id = ?", collectionId).
		Find((&resultCollection)).Error

	if err != nil {
		intServerError(w, err)
		return
	}

	libOwn := database.LibraryData{}

	err = database.DB.Where("user_id = ?", auth.ID).Find(&libOwn).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	if libOwn.ID != resultCollection.LibraryID {
		unauthorizedRequest(w, errors.New("user not allowed"))
		return
	}

	err = database.DB.Model(&resultCollection).
		Updates(database.LibraryCollection{
			Status:       e.Status,
			Availability: e.Availability}).
		Error

	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Collection status changed",
	}.responseFormatter(w)
}
