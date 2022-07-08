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
		AllowedHeaders: []string{"Authorization", "Content-Type", "Account-auth"},
	})

	r := mux.NewRouter()

	// NOTE: db handler
	r.HandleFunc("/db/migrate", testMigrate).Methods("GET")
	r.HandleFunc("/db/seed/profile", testSeedProfile).Methods("GET")
	r.HandleFunc("/db/seed/book", testSeedBook).Methods("GET")
	r.HandleFunc("/db/seed/library", testSeedLibrary).Methods("GET")

	r.HandleFunc("/profile", profileInformation).Methods("GET")
	r.HandleFunc("/profile/borrow", profileBorrow).Methods("GET")
	r.HandleFunc("/profile/login", loginForm).Methods("POST")
	r.HandleFunc("/profile/login/google", loginGoogle).Methods("POST")
	r.HandleFunc("/profile/update/password", updatePassword).Methods("POST")
	r.HandleFunc("/profile/update/email", updateEmail).Methods("POST")
	r.HandleFunc("/profile/update/picture", updatePicture).Methods("POST")
	r.HandleFunc("/profile/update/username", updateUsername).Methods("POST")
	r.HandleFunc("/profile/validate/resend", resendEmail).Methods("GET")
	r.HandleFunc("/profile/permission", profileAccessPermission).Methods("GET")
	r.HandleFunc("/profile/validate", emailValidate).Methods("GET")

	r.HandleFunc("/book", allBooks).Methods("GET")
	r.HandleFunc("/book/{slug}", singleBook).Methods("GET")

	r.HandleFunc("/library", allLibraries).Methods("GET")
	r.HandleFunc("/library/{id}", singleLibrary).Methods("GET")

	r.HandleFunc("/paper", allPapers).Methods("GET")
	r.HandleFunc("/paper/{id}", singlePaper).Methods("GET")

	http.Handle("/", r)
	utility.InfoPrint(1, fmt.Sprintf("service at port %s", port))

	http.ListenAndServe(":"+port, c.Handler(r))
}
