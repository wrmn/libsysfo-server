package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	r := mux.NewRouter()
	// r.HandleFunc("/", LibrariesHandler).Methods("GET")
	r.HandleFunc("/library", LibrariesHandler).Methods("GET")
	r.HandleFunc("/library/{id}", LibraryHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("Listing for" + port)

	fmt.Println(http.ListenAndServe(":"+port, r))
}

func LibrariesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(msg)
}

func LibraryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id > 11 {
		json.NewEncoder(w).Encode("out of range or invalid id")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(msg.Data[id-1])
}
