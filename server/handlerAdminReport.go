package server

import (
	"libsysfo-server/database"
	"libsysfo-server/utility"
	"libsysfo-server/utility/report"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func booksReport(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	q := r.URL.Query()
	category := getCatgory(q)
	timeRange, invalid := getTimeRange(q, w)
	if invalid {
		return
	}

	bookQuery := []database.LibraryCollection{}

	db := database.DB.
		Where("library_id = ? AND created_at BETWEEN ? AND ?",
			libraryData.ID,
			timeRange[0],
			timeRange[1]).Preload("Book", func(db *gorm.DB) *gorm.DB {
		return database.DB.Preload("BookDetail")
	}).Order("id desc").Find(&bookQuery)
	if invalid := databaseException(w, db); invalid {
		return
	}

	s := "Sheet1"
	xlsx := report.Table{
		Header: []report.MainHeader{{
			Name:  "Created At",
			Value: time.Now().Format("2 Jan 2006 15:04:05"),
		}, {
			Name:  "Title",
			Value: "Books Report",
		}},
		Table: report.BookReport}

	f := xlsx.CreateMainTable(s)

	i := 0
	for _, d := range bookQuery {
		if utility.Compare(category, d.Book.BookDetail.Category) || len(category) == 0 {
			i++
			col := i + len(xlsx.Header) + 1
			appendBookRequest(f, s, col, i, d)
		}
	}

	xlsx.Data = i
	if i == 0 {
		badRequest(w, "no data to be reported")
		return
	}

	xlsx.Styling(s, f)

	body, err := f.WriteToBuffer()
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Write(body.Bytes())
}

func bookReport(w http.ResponseWriter, r *http.Request) {
	_, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	collectionId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "invalid id request")
		return
	}

	resultCollection, invalid := findCollectionById(collectionId, w)
	if invalid {
		return
	}

	s := "Sheet1"
	xlsx := report.Table{
		Header: generateBookHeader(resultCollection),
		Table:  report.BorrowReport,
		Data:   len(resultCollection.Borrow),
	}

	f := xlsx.CreateMainTable(s)

	for i, d := range resultCollection.Borrow {
		col := i + 2 + len(xlsx.Header)
		appendBorrowReport(f, s, col, i+1, d)
	}

	xlsx.Styling(s, f)

	body, err := f.WriteToBuffer()
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Write(body.Bytes())
}

func borrowReport(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	q := r.URL.Query()
	timeRange, invalid := getTimeRange(q, w)
	if invalid {
		return
	}

	collectionsData := []database.LibraryCollection{}
	err := database.DB.Where("library_id = ?", libraryData.ID).
		Preload("Borrow", func(db *gorm.DB) *gorm.DB {
			return database.DB.Where("created_at BETWEEN ? AND ?", timeRange[0], timeRange[1]).
				Preload("Collection.Book").
				Preload("User.ProfileData")
		}).Find(&collectionsData).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	s := "Sheet1"
	xlsx := report.Table{
		Header: []report.MainHeader{{
			Name:  "Created At",
			Value: time.Now().Format("2 Jan 2006 15:04:05"),
		}, {
			Name:  "Title",
			Value: "Books Report",
		}},
		Table: report.BorrowReport}

	f := xlsx.CreateMainTable(s)
	i := 0
	for _, k := range collectionsData {
		for _, l := range k.Borrow {
			i++
			col := i + len(xlsx.Header) + 1
			appendBorrowReport(f, s, col, i+1, l)
		}
	}

	xlsx.Data = i
	if i == 0 {
		badRequest(w, "no data to be reported")
		return
	}

	xlsx.Styling(s, f)

	body, err := f.WriteToBuffer()
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Write(body.Bytes())
}
