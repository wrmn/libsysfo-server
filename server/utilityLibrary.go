package server

import (
	"fmt"
	"libsysfo-server/database"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func getLibraryBook(libId int, r *http.Request) (bookData []bookResponse, err error) {
	bookQuery := []database.LibraryCollection{}
	q := r.URL.Query()
	var db *gorm.DB

	if q.Has("sn") {
		sn := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("sn")))
		db = database.DB.Where("library_id = ? AND LOWER(serial_number) LIKE ?", libId, sn)

	} else {
		db = database.DB.Where("library_id = ?", libId)
	}

	err = db.
		Preload("Book", func(db *gorm.DB) *gorm.DB {
			return database.DB.Preload("BookDetail")
		}).Order("id desc").Find(&bookQuery).Error

	if err != nil {
		return
	}

	for _, d := range bookQuery {
		bookData = append(bookData, bookResponse{
			Id:          d.ID,
			Title:       d.Book.Title,
			Image:       d.Book.Image,
			Author:      d.Book.Author,
			Slug:        d.Book.Slug,
			ReleaseDate: d.Book.BookDetail.ReleaseDate,
			Language:    d.Book.BookDetail.Language,
			Description: d.Book.BookDetail.Description,
			Country:     d.Book.BookDetail.Country,
			Publisher:   d.Book.BookDetail.Publisher,
			PageCount:   d.Book.BookDetail.PageCount,
			Category:    d.Book.BookDetail.Category,
			Status: &libraryCollectionResponse{
				Id:           d.ID,
				SerialNumber: d.SerialNumber,
				Availability: d.Availability,
				Status:       d.Status,
			},
		})
	}

	return
}

func getLibraryPaper(libId int) (paperData []paperResponse, err error) {
	paperQuery := []database.LibraryPaper{}
	err = database.DB.Where("library_id = ?", libId).Find(&paperQuery).Error
	if err != nil {
		return
	}

	for _, c := range paperQuery {
		paperData = append(paperData, paperResponse{
			Id:          c.ID,
			Title:       c.Title,
			Subject:     c.Subject,
			Abstract:    c.Abstract,
			Type:        c.Type,
			Description: c.Description,
			Access:      c.Access,
		})
	}
	return
}

func setBookResponse(result database.Book) (bookRespBody bookResponse) {
	bookRespBody.Title = result.Title
	bookRespBody.Image = result.Image
	bookRespBody.Author = result.Author
	bookRespBody.Slug = result.Slug
	bookRespBody.Source = result.Source
	bookRespBody.ReleaseDate = result.BookDetail.ReleaseDate
	bookRespBody.Description = result.BookDetail.Description
	bookRespBody.Language = result.BookDetail.Language
	bookRespBody.Country = result.BookDetail.Country
	bookRespBody.PageCount = int(result.BookDetail.PageCount)
	bookRespBody.Publisher = result.BookDetail.Publisher
	bookRespBody.Category = result.BookDetail.Category
	return
}

func paginator(r *http.Request, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func bookFilter(r *http.Request, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
		if q.Has("keyword") {
			keyword := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("keyword")))
			db = db.Where("LOWER(title) like ?", keyword).
				Or("LOWER(author) like ?", keyword)
		}
		return db
	}
}

func bookDetailFilter(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()

		if q.Has("language") {
			if q.Has("category") {
				category := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("category")))
				db = db.Where("LOWER(language) = ? AND LOWER(category) = ?", strings.ToLower(q.Get("language")), category)
				return db
			}
			db = db.Where("LOWER(language) = ? ", strings.ToLower(q.Get("language")))
			return db
		} else if q.Has("category") {
			category := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("category")))
			db = db.Where("LOWER(category) like ? ", category)
			return db
		}
		return db
	}
}

func paperFilter(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		if q.Has("keyword") {
			keyword := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("keyword")))
			db = db.Where("LOWER(title) like ?", keyword).
				Or("array_to_string(subject, ',', ' ') like ?", keyword).
				Or("description::TEXT like ?", keyword)
		}
		return db
	}
}

func (data paginate) generate(r *http.Request, page int) (result paginate) {
	// NOTE:change to https when deployed
	link := "http://" + r.Host + r.URL.Path + "?page="
	if page > 1 {
		result.Prev = link + strconv.Itoa(page-1)
	}
	result.Current = link + strconv.Itoa(page)
	result.Data = data.Data
	if result.Data != 0 {
		result.Next = link + strconv.Itoa(page+1)
	}
	return
}

func setStatus(d database.LibraryCollectionBorrow) string {
	if d.ReturnedAt != nil {
		return "finished"
	}
	if d.TakedAt != nil {
		return "taked"
	}
	if d.CanceledAt != nil {
		return "canceled"
	}
	if d.AcceptedAt != nil {
		return "accepted"
	}
	return "requested"
}

func findCollectionById(collectionId int) (result database.LibraryCollection, err error) {
	err = database.DB.Where("id = ?", collectionId).Preload("Book.BookDetail").
		Find((&result)).Error
	return
}
