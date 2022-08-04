package server

import (
	"libsysfo-server/database"
	"libsysfo-server/utility"
	"libsysfo-server/utility/report"
	"net/http"
	"net/url"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func getTimeRange(q url.Values, w http.ResponseWriter) ([]string, bool) {
	timeRange := []string{"20060101", time.Now().Format("20060102")}
	if q.Has("from") {
		t, err := time.Parse("20060102", q.Get("from"))
		if err != nil {
			badRequest(w, "invalid parameter")
			return nil, true
		}
		timeRange[0] = t.Format("20060102")
	}
	if q.Has("to") {
		t, err := time.Parse("20060102", q.Get("to"))
		if err != nil {
			badRequest(w, "invalid parameter")
			return nil, true
		}
		timeRange[1] = t.Add(time.Hour * 24).Format("20060102")
	}
	return timeRange, false
}

func getCatgory(q url.Values) (category []string) {
	for key, values := range q {
		if key == "category" {
			category = values
		}
	}
	return category
}

func generateBookHeader(d database.LibraryCollection) []report.MainHeader {
	return []report.MainHeader{{
		Name:  "Report Created At",
		Value: time.Now().Format("2 Jan 2006 15:04:05"),
	}, {
		Name:  "Collection Created At",
		Value: d.CreatedAt.Format("2 Jan 2006 15:04:05"),
	}, {
		Name:  "Book Title",
		Value: d.Book.Title,
	}, {
		Name:  "Serial Number",
		Value: d.SerialNumber,
	}, {
		Name:  "Category",
		Value: d.Book.BookDetail.Category,
	}, {
		Name:  "Author",
		Value: d.Book.Author,
	}, {
		Name:  "Release Date",
		Value: d.Book.BookDetail.ReleaseDate,
	}, {
		Name:  "Publisher",
		Value: d.Book.BookDetail.Publisher,
	}, {
		Name:  "Language",
		Value: d.Book.BookDetail.Language,
	}, {
		Name:  "Country",
		Value: d.Book.BookDetail.Country,
	}, {
		Name:  "Page",
		Value: d.Book.BookDetail.PageCount,
	}, {
		Name:  "Status",
		Value: d.Status,
	}, {
		Name:  "Borrow Total",
		Value: len(d.Borrow),
	}}
}

func appendBookRequest(f *excelize.File, s string, col int, i int, d database.LibraryCollection) {
	var count int64
	database.DB.
		Model(&database.LibraryCollectionBorrow{}).
		Where("collection_id = ?", d.ID).
		Count(&count)
	f.SetCellValue(s, report.Cell("A", col), i)
	f.SetCellValue(s, report.Cell("B", col), d.CreatedAt.Format("2 Jan 2006 15:04:05"))
	f.SetCellValue(s, report.Cell("C", col), d.SerialNumber)
	f.SetCellValue(s, report.Cell("D", col), d.Book.Title)
	f.SetCellValue(s, report.Cell("E", col), d.Book.BookDetail.Category)
	f.SetCellValue(s, report.Cell("F", col), d.Book.Author)
	f.SetCellValue(s, report.Cell("G", col), d.Book.BookDetail.ReleaseDate)
	f.SetCellValue(s, report.Cell("H", col), d.Book.BookDetail.Publisher)
	f.SetCellValue(s, report.Cell("I", col), d.Book.BookDetail.Language)
	f.SetCellValue(s, report.Cell("J", col), d.Book.BookDetail.Country)
	f.SetCellValue(s, report.Cell("K", col), d.Book.BookDetail.PageCount)
	f.SetCellValue(s, report.Cell("L", col), utility.AvailabilityString(d.Availability))
	f.SetCellValue(s, report.Cell("M", col), utility.StatusString(d.Status))
	f.SetCellValue(s, report.Cell("N", col), count)
}

func appendBorrowReport(f *excelize.File, s string, col int, i int, d database.LibraryCollectionBorrow) {
	f.SetCellValue(s, report.Cell("A", col), i)
	f.SetCellValue(s, report.Cell("B", col), setBorrowStatus(d))
	f.SetCellValue(s, report.Cell("C", col), d.CreatedAt.Format("2 Jan 2006 15:04:05"))
	f.SetCellValue(s, report.Cell("D", col), getDate(d.AcceptedAt))
	f.SetCellValue(s, report.Cell("E", col), getDate(d.TakedAt))
	f.SetCellValue(s, report.Cell("F", col), getDate(d.ReturnedAt))
	f.SetCellValue(s, report.Cell("G", col), getDate(d.CanceledAt))
	f.SetCellValue(s, report.Cell("H", col), d.Collection.Book.Title)
	f.SetCellValue(s, report.Cell("I", col), d.Collection.SerialNumber)
	f.SetCellValue(s, report.Cell("J", col), d.User.ProfileData.Name)
	f.SetCellValue(s, report.Cell("K", col), getString(d.User.Username))
	f.SetCellValue(s, report.Cell("L", col), d.User.Email)
}
