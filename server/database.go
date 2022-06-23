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
