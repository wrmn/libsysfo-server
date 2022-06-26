package server

import (
	"fmt"
	"libsysfo-server/utility"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/booktest", booksTestHandler).Methods("GET")
	r.HandleFunc("/db/migrate", testMigrate).Methods("GET")
	r.HandleFunc("/db/seed/profile", testSeedProfile).Methods("GET")
	r.HandleFunc("/db/seed/book", testSeedBook).Methods("GET")
	r.HandleFunc("/db/seed/library", testSeedLibrary).Methods("GET")
	http.Handle("/", r)
	utility.InfoPrint(1, fmt.Sprintf("service at port %s", port))
	http.ListenAndServe(":"+port, r)
}
