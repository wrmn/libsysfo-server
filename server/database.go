package server

import (
	"encoding/json"
	"libsysfo-server/database"
	"net/http"
)

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
