package server

import (
	"fmt"
	"libsysfo-server/database"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func getLibraryPaper(libId int, r *http.Request) (paperData []paperResponse, err error) {
	paperQuery := []database.LibraryPaper{}
	q := r.URL.Query()
	var db *gorm.DB
	err = database.DB.Where("library_id = ?", libId).Find(&paperQuery).Error
	if err != nil {
		return
	}

	if q.Has("pkw") {
		pkw := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("pkw")))
		db = database.DB.Where("library_id = ? AND LOWER(title) LIKE ?", libId, pkw)
	} else {
		db = database.DB.Where("library_id = ?", libId)
	}

	err = db.Order("id desc").Find(&paperQuery).Error
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

func setBookResponse(result database.Book) bookResponse {
	return bookResponse{
		Title:       result.Title,
		Image:       result.Image,
		Author:      result.Author,
		Slug:        result.Slug,
		Source:      result.Source,
		ReleaseDate: result.BookDetail.ReleaseDate,
		Description: result.BookDetail.Description,
		Language:    result.BookDetail.Language,
		Country:     result.BookDetail.Country,
		PageCount:   int(result.BookDetail.PageCount),
		Publisher:   result.BookDetail.Publisher,
		Category:    result.BookDetail.Category,
	}
}

func setPaperResponse(result database.LibraryPaper) paperResponse {
	return paperResponse{
		Id:          result.ID,
		Title:       result.Title,
		Subject:     result.Subject,
		Abstract:    result.Abstract,
		Type:        result.Type,
		Description: result.Description,
		Access:      result.Access,
	}
}

func appendBorrowData(data []database.LibraryCollectionBorrow) (respBody []profileCollectionBorrowResponse) {
	for _, d := range data {
		respBody = append(respBody, formatBorrowData(d))
	}
	return
}

func appendPermissionData(data []database.LibraryPaperPermission) (respBody []profilePermissionResponse) {
	for _, j := range data {
		respBody = append(respBody, formatPermissionData(j))
	}
	return
}

func formatBorrowData(d database.LibraryCollectionBorrow) profileCollectionBorrowResponse {
	return profileCollectionBorrowResponse{
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
		Status:       setBorrowStatus(d),
	}
}

func formatPermissionData(j database.LibraryPaperPermission) profilePermissionResponse {
	status := setPermissionStatus(j)
	return profilePermissionResponse{
		CreatedAt:        j.CreatedAt,
		AcceptedAt:       j.AcceptedAt,
		CanceledAt:       j.CanceledAt,
		Id:               j.ID,
		PaperId:          j.PaperID,
		PaperTitle:       j.Paper.Title,
		PaperSubject:     j.Paper.Subject,
		PaperDescription: j.Paper.Description,
		PaperType:        j.Paper.Type,
		LibraryId:        j.Paper.LibraryID,
		Library:          j.Paper.Library.Name,
		Purpose:          j.Purpose,
		UserId:           j.UserID,
		UserName:         j.User.ProfileData.Name,
		Status:           status,
	}
}

func appendAccessData(data []database.LibraryPaperAccess) (resp []time.Time) {
	for _, k := range data {
		resp = append(resp, k.CreatedAt)
	}
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

func setBorrowStatus(d database.LibraryCollectionBorrow) string {
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

func getDate(t *time.Time) string {
	if t != nil {
		return t.Format("2 Jan 2006 15:04:05")
	}
	return "-"
}

func getString(s *string) string {
	if s != nil {
		return *s
	}
	return "-"
}

func getAccessValue(d bool) string {
	if d {
		return "Accessible"
	}
	return "Inaccessible"
}

func setPermissionStatus(p database.LibraryPaperPermission) string {
	if p.CanceledAt != nil {
		return "canceled"
	}
	if p.AcceptedAt != nil {
		return "accepted"
	}
	return "requested"
}

func findCollectionById(collectionId int, w http.ResponseWriter) (result database.LibraryCollection, invalid bool) {
	db := database.DB.Where("id = ?", collectionId).Preload("Book.BookDetail").
		Preload("Borrow.Collection.Book").
		Preload("Borrow.User.ProfileData").
		Find((&result))

	invalid = databaseException(w, db)
	return
}

func cancelOtherBorrow(currentData database.LibraryCollectionBorrow) (err error) {
	data := []database.LibraryCollectionBorrow{}
	err = database.DB.
		Where("NOT id = ? AND collection_id = ? AND returned_at IS NULL",
			currentData.ID,
			currentData.CollectionID,
		).
		Find(&data).
		Error
	now := time.Now()

	for i := range data {
		data[i].CanceledAt = &now
		database.DB.Save(&data[i])
	}
	return
}

func generateCollectionResponse(data database.LibraryCollection) libraryCollectionResponse {
	return libraryCollectionResponse{
		Id:           data.ID,
		SerialNumber: data.SerialNumber,
		Availability: data.Availability,
		Status:       data.Status,
	}
}
