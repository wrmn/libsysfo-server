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
	// remove if unused
	// r.HandleFunc("/db/migrate", testMigrate).Methods("GET")
	// r.HandleFunc("/db/seed/profile", testSeedProfile).Methods("GET")
	// r.HandleFunc("/db/seed/book", testSeedBook).Methods("GET")
	// r.HandleFunc("/db/seed/library", testSeedLibrary).Methods("GET")

	r.HandleFunc("/login", superAdminLogin).Methods("POST")
	r.HandleFunc("/feedback", getFeedback).Methods("GET")
	r.HandleFunc("/library/new", newLibrary).Methods("POST")
	r.HandleFunc("/third-party", getThirdParty).Methods("GET")

	r.HandleFunc("/feedback", newFeedback).Methods("POST")
	r.HandleFunc("/notification", getNotification).Methods("GET")
	r.HandleFunc("/notification", readNotification).Methods("POST")

	r.HandleFunc("/admin/library", adminInformation).Methods("GET")
	r.HandleFunc("/admin/library/login", adminLogin).Methods("POST")
	r.HandleFunc("/admin/library/dashboard", libraryDashboard).Methods("GET")

	r.HandleFunc("/admin/library/collection", libraryCollections).Methods("GET")
	r.HandleFunc("/admin/library/collection/new", libraryAddCollection).Methods("POST")
	r.HandleFunc("/admin/library/collection/{id}", librarySingleCollection).Methods("GET")
	r.HandleFunc("/admin/library/collection/{id}/update", libraryUpdateCollection).Methods("POST")

	r.HandleFunc("/admin/library/borrow", libraryBorrow).Methods("GET")
	r.HandleFunc("/admin/library/borrow/find", findBorrow).Methods("GET")
	r.HandleFunc("/admin/library/borrow/action", actionBorrow).Methods("POST")

	r.HandleFunc("/admin/library/paper", libraryPaper).Methods("GET")
	r.HandleFunc("/admin/library/paper/new", libraryAddPaper).Methods("POST")
	r.HandleFunc("/admin/library/paper/{id}", librarySinglePaper).Methods("GET")
	r.HandleFunc("/admin/library/paper/{id}/update", libraryUpdatePaper).Methods("POST")
	r.HandleFunc("/admin/library/paper/{id}/update/file", libraryUpdatePaperFile).Methods("POST")

	r.HandleFunc("/admin/library/permission", libraryPermission).Methods("GET")
	r.HandleFunc("/admin/library/permission/find", findPermission).Methods("GET")
	r.HandleFunc("/admin/library/permission/action", actionPermission).Methods("POST")
	r.HandleFunc("/admin/library/permission/{id}/access", accessHistory).Methods("GET")

	r.HandleFunc("/admin/library/report/book", booksReport).Methods("GET")
	r.HandleFunc("/admin/library/report/book/{id}", bookReport).Methods("GET")
	r.HandleFunc("/admin/library/report/borrow", borrowReport).Methods("GET")
	r.HandleFunc("/admin/library/report/paper", papersReport).Methods("GET")
	r.HandleFunc("/admin/library/report/paper/{id}", paperReport).Methods("GET")
	r.HandleFunc("/admin/library/report/permission", permissionReport).Methods("GET")

	r.HandleFunc("/admin/library/user/find", libraryUserFind).Methods("GET")
	r.HandleFunc("/admin/library/user/{id}", libraryUser).Methods("GET")
	r.HandleFunc("/admin/library/user/{id}/borrow", libraryUserBorrow).Methods("GET")
	r.HandleFunc("/admin/library/user/{id}/permission", libraryUserPermission).Methods("GET")

	r.HandleFunc("/admin/library/update", libraryGeneral).Methods("POST")
	r.HandleFunc("/admin/library/update/image", libraryImage).Methods("POST")
	r.HandleFunc("/admin/library/update/image/content", libraryContentImage).Methods("POST")
	r.HandleFunc("/admin/library/update/location", libraryLocation).Methods("POST")

	r.HandleFunc("/profile", profileInformation).Methods("GET")
	r.HandleFunc("/profile/login", loginForm).Methods("POST")
	r.HandleFunc("/profile/login/google", loginGoogle).Methods("POST")
	r.HandleFunc("/profile/register", registerForm).Methods("POST")

	r.HandleFunc("/profile/update/password", updatePassword).Methods("POST")
	r.HandleFunc("/profile/update/email", updateEmail).Methods("POST")
	r.HandleFunc("/profile/update/picture", updatePicture).Methods("POST")
	r.HandleFunc("/profile/update/profile", updateProfile).Methods("POST")
	r.HandleFunc("/profile/update/username", updateUsername).Methods("POST")

	r.HandleFunc("/profile/permission", profileAccessPermission).Methods("GET")
	r.HandleFunc("/profile/permission/new", profileNewPermission).Methods("POST")
	r.HandleFunc("/profile/permission/read/{id}", profileReadPaper).Methods("GET")

	r.HandleFunc("/profile/validate", emailValidate).Methods("GET")
	r.HandleFunc("/profile/validate/resend", resendEmail).Methods("GET")

	r.HandleFunc("/profile/borrow", profileBorrow).Methods("GET")
	r.HandleFunc("/profile/borrow/new", borrowNewBook).Methods("POST")

	r.HandleFunc("/book", allBooks).Methods("GET")
	r.HandleFunc("/book/{slug}", singleBook).Methods("GET")

	r.HandleFunc("/library", allLibraries).Methods("GET")
	r.HandleFunc("/library/{id}", singleLibrary).Methods("GET")

	r.HandleFunc("/paper", allPapers).Methods("GET")
	r.HandleFunc("/paper/{id}", singlePaper).Methods("GET")

	r.MethodNotAllowedHandler = handleNotAllowed()
	r.NotFoundHandler = handleNotFound()

	http.Handle("/", r)
	utility.InfoPrint(1, fmt.Sprintf("service at port %s", port))

	http.ListenAndServe(":"+port, c.Handler(r))
}
