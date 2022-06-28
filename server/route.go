package server

import (
	"fmt"
	"libsysfo-server/utility"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Serve(port string) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	r := mux.NewRouter()
	// NOTE: testing endpoint
	r.HandleFunc("/booktest", booksTestHandler).Methods("GET")
	r.HandleFunc("/db/migrate", testMigrate).Methods("GET")
	r.HandleFunc("/db/seed/profile", testSeedProfile).Methods("GET")
	r.HandleFunc("/db/seed/book", testSeedBook).Methods("GET")
	r.HandleFunc("/db/seed/library", testSeedLibrary).Methods("GET")

	r.HandleFunc("/profile/login/google", loginGoogle).Methods("POST")
	r.HandleFunc("/profile/validate", emailValidate).Methods("GET")
	http.Handle("/", r)
	utility.InfoPrint(1, fmt.Sprintf("service at port %s", port))
	http.ListenAndServe(":"+port, c.Handler(r))
}
