package server

import (
	"encoding/json"
	"fmt"
	"libsysfo-server/database"
	"libsysfo-server/utility"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func Serve(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/library", librariesHandler).Methods("GET")
	r.HandleFunc("/library/{id}", libraryHandler).Methods("GET")
	r.HandleFunc("/book", booksHandler).Methods("GET")
	r.HandleFunc("/db/test", testDb).Methods("GET")

	http.Handle("/", r)

	utility.InfoPrint(1, fmt.Sprintf("service at port %s", port))
	http.ListenAndServe(":"+port, r)
}

type User struct {
	Name    string
	NewName string
}

func testDb(w http.ResponseWriter, r *http.Request) {

	err := database.DB.Migrator().CreateTable(&User{})

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode("database terbaca jing")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode("test is successful")
}

func librariesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(librariesData)
}

func libraryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id > 11 {
		json.NewEncoder(w).Encode("out of range or invalid id")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(librariesData.Data[id-1])
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(booksData)
}
