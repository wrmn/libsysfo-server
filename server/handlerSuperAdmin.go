package server

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"libsysfo-server/database"
	"libsysfo-server/utility/email"
	"net/http"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/lib/pq"
)

func superAdminLogin(w http.ResponseWriter, r *http.Request) {
	user, err := getLoginData(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	} else if user.AccountType != 1 {
		badRequest(w, "user not allowed")
		return
	}
	loginHandler(w, user)
}

func getFeedback(w http.ResponseWriter, r *http.Request) {
	_, invalid := checkIfAllowed(1, w, r)
	if invalid {
		return
	}
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page == 0 {
		page = 1
	}

	feedback := []database.Feedback{}
	db := database.DB.Scopes(paginator(r, 10)).Order("created_at DESC").Find(&feedback)
	if invalid := databaseException(w, db); invalid {
		return
	}

	paginateData := paginate{Data: len(feedback)}.generate(r, page)
	response{
		Data: responseBody{
			Feedback: feedback,
			Paginate: &paginateData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func newLibrary(w http.ResponseWriter, r *http.Request) {
	if _, invalid := checkIfAllowed(1, w, r); invalid {
		return
	}

	lastAcc := database.ProfileAccount{}
	database.DB.Order("id desc").First(&lastAcc)

	var e newLibraryRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()
	err := decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}
	password := gofakeit.Gamertag()
	uname := gofakeit.Gamertag()
	user := database.ProfileAccount{
		ID:          lastAcc.ID + 1,
		Username:    &uname,
		Email:       e.Email,
		AccountType: 2,
		Password:    fmt.Sprintf("%x", md5.Sum([]byte(password))),
	}
	err = database.DB.Create(&user).Error
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	content := fmt.Sprintf("<html><head></head><body><p>Hello,</p>Perpustakaan anda telah terdaftar, silahkan login di library sistem information for admin dengan email anda atau username %s dan password %s dan lengkapi informasi perpustakaan anda</body>	</html>",
		uname, password,
	)

	emailBody := email.Content{
		Subject:     "Pendaftaran perpustakaan",
		HtmlContent: content,
	}

	receiver := email.ToData{
		Name:  e.Name,
		Email: e.Email,
	}

	err = emailBody.SendEmail(receiver)
	if err != nil {
		fmt.Println(err.Error())
	}

	lastLib := database.LibraryData{}
	database.DB.Order("id desc").First(&lastLib)
	library := database.LibraryData{
		ID:            lastLib.ID + 1,
		UserID:        user.ID,
		Name:          e.LibraryName,
		Address:       "",
		Coordinate:    pq.Float64Array{100.36262663239825, -0.9225479705730635},
		Description:   "",
		ImagesMain:    "https://upload.wikimedia.org/wikipedia/commons/2/21/Biblioth%C3%A8que_de_l%27Assembl%C3%A9e_Nationale_%28Lunon%29.jpg",
		ImagesContent: pq.StringArray{},
		Webpage:       "",
	}
	err = database.DB.Create(&library).Error
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func getThirdParty(w http.ResponseWriter, r *http.Request) {
	_, invalid := checkIfAllowed(1, w, r)
	if invalid {
		return
	}
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page == 0 {
		page = 1
	}

	jobs := []database.ThirdPartyJobs{}
	db := database.DB.Scopes(paginator(r, 10)).Order("created_at DESC").Find(&jobs)
	if invalid := databaseException(w, db); invalid {
		return
	}

	paginateData := paginate{Data: len(jobs)}.generate(r, page)
	response{
		Data: responseBody{
			Jobs:     jobs,
			Paginate: &paginateData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}
