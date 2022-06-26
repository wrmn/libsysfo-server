package server

import (
	"encoding/json"
	"libsysfo-server/database"
	"net/http"
)

func testMigrate(w http.ResponseWriter, r *http.Request) {
	err := database.DB.AutoMigrate(
		&database.Book{},
		&database.BookDetail{},
		&database.LibraryData{},
		&database.LibraryCollection{},
		&database.LibraryCollectionBorrow{},
		&database.LibraryPaper{},
		&database.LibraryPaperPermission{},
		&database.LibraryPaperAccess{},
		&database.ProfileAccount{},
		&database.ProfileData{},
	)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("table created")

}

func testSeedProfile(w http.ResponseWriter, r *http.Request) {
	database.SeedProfileAccount()
	database.SeedProfileData()
}

func testSeedBook(w http.ResponseWriter, r *http.Request) {
	// NOTE: this need API make sure seed it when it needed to be
	database.SeedBook()
	database.SeedBookDetail()
}

func testSeedLibrary(w http.ResponseWriter, r *http.Request) {
	database.SeedLibraryData()
	database.SeedLibraryCollection()
	database.SeedLibraryPaper()
	database.SeedLibraryPaperPermission()
}
