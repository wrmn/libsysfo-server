package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type ImagesData struct {
	Main    string   `json:"main"`
	Content []string `json:"content"`
}

type NameData struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	Coordinate  []float64  `json:"coordinate"`
	Description string     `json:"description"`
	Images      ImagesData `json:"images"`
}

type NameDatas struct {
	Data []NameData `json:"data"`
}

func main() {

	port := os.Getenv("PORT")
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("Listing for" + port)

	fmt.Println(http.ListenAndServe(":"+port, r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	msg := "Semangat ayang, rang baru deploy, wkwkw. ko baru bisa masuk internet, program yang di deploy alun baubah lai"

	time.Sleep(1 * time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(msg)
}
