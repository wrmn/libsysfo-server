package server

import (
	"context"
	"fmt"
	"libsysfo-server/utility"
	"net/http"
	"os"

	"github.com/codedius/imagekit-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Serve(port string) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	r := mux.NewRouter()
	// NOTE: testing endpoint
	r.HandleFunc("/img/post", imgKitUpload).Methods("POST")

	// NOTE: db handler
	r.HandleFunc("/db/migrate", testMigrate).Methods("GET")
	r.HandleFunc("/db/seed/profile", testSeedProfile).Methods("GET")
	r.HandleFunc("/db/seed/book", testSeedBook).Methods("GET")
	r.HandleFunc("/db/seed/library", testSeedLibrary).Methods("GET")

	r.HandleFunc("/profile/login", loginForm).Methods("POST")
	r.HandleFunc("/profile/login/google", loginGoogle).Methods("POST")
	r.HandleFunc("/profile/validate", emailValidate).Methods("GET")
	r.HandleFunc("/profile", profileInformation).Methods("GET")
	r.HandleFunc("/profile/borrow", profileBorrow).Methods("GET")
	r.HandleFunc("/profile/permission", profileAccessPermission).Methods("GET")

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

func imgKitUpload(w http.ResponseWriter, r *http.Request) {
	opts := imagekit.Options{
		PublicKey:  os.Getenv("IMAGEKIT_PUBLIC_KEY"),
		PrivateKey: os.Getenv("IMAGEKIT_PRIVATE_KEY"),
	}

	ik, err := imagekit.NewClient(&opts)
	if err != nil {
		intServerError(w, err)
		return
	}
	ur := imagekit.UploadRequest{
		File:              "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b6/Image_created_with_a_mobile_phone.png/800px-Image_created_with_a_mobile_phone.png",
		FileName:          "testing",
		UseUniqueFileName: false,
		Tags:              []string{"testing", "test"},
		Folder:            "/",
		IsPrivateFile:     false,
		CustomCoordinates: "",
		ResponseFields:    nil,
	}

	ctx := context.Background()

	upr, err := ik.Upload.ServerUpload(ctx, &ur)
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Data: responseBody{
			Profile: &upr,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)

}
