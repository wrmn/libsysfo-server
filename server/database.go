package server

import (
	"encoding/json"
	"libsysfo-server/database"
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func testExcel(w http.ResponseWriter, r *http.Request) {
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "B2", 100)
	f.SetCellValue("Sheet1", "A1", 50)

	now := time.Now()

	f.SetCellValue("Sheet1", "A4", now.Format(time.ANSIC))
	body, err := f.WriteToBuffer()
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Write(body.Bytes())
}

func testMigrate(w http.ResponseWriter, _ *http.Request) {
	err := database.DB.AutoMigrate(
		&database.Book{},
		&database.BookDetail{},
		&database.ProfileAccount{},
		&database.ProfileData{},
		&database.LibraryData{},
		&database.LibraryCollection{},
		&database.LibraryCollectionBorrow{},
		&database.LibraryPaper{},
		&database.LibraryPaperPermission{},
		&database.LibraryPaperAccess{},
		&database.Feedback{},
		&database.ThirdPartyJobs{},
	)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("table created")

}

func testSeedProfile(_ http.ResponseWriter, _ *http.Request) {
	database.SeedProfileAccount()
	database.SeedProfileData()
}

func testSeedBook(_ http.ResponseWriter, _ *http.Request) {
	database.SeedBookLocal()
	// NOTE: this need API make sure seed it when it needed to be
	database.SeedBook()
	database.SeedBookDetail()
}

func testSeedLibrary(_ http.ResponseWriter, _ *http.Request) {
	database.SeedLibraryData()
	database.SeedLibraryCollection()
	database.SeedLibraryCollectionBorrow()
	database.SeedLibraryPaper()
	database.SeedLibraryPaperPermission()
	database.SeedLibraryPaperAccess()
}
