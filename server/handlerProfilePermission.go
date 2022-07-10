package server

import (
	"errors"
	"io/ioutil"
	"libsysfo-server/database"
	"libsysfo-server/utility/cred"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func profileAccessPermission(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	data := database.ProfileAccount{}
	cred := tokenData.Claims.(*cred.TokenModel)
	database.DB.Where("email = ?", cred.Email).Or("username = ?", cred.Username).
		Preload("ProfileData").First(&data)

	permissionData := searchPermission(data.ID)

	response{
		Data: responseBody{
			Permission: permissionData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func profileReadPaper(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authVerification(r)
	id := mux.Vars(r)["id"]
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	data := database.ProfileAccount{}
	cred := tokenData.Claims.(*cred.TokenModel)
	database.DB.Where("email = ?", cred.Email).Or("username = ?", cred.Username).
		Preload("ProfileData").First(&data)

	permissionData := database.LibraryPaperPermission{}
	row := database.DB.
		Where("id = ? AND user_id = ?", id, data.ID).
		Preload("Paper").Find(&permissionData).RowsAffected

	if row == 0 {
		unauthorizedRequest(w, errors.New("data not found"))
	}

	resp, err := http.Get(permissionData.Paper.PaperUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	access := database.LibraryPaperAccess{
		PermissionID: permissionData.ID,
		CreatedAt:    time.Now(),
	}

	database.DB.Save(&access)

	w.Write(body)
}

func searchPermission(id int) (respBody []profilePermissionResponse) {
	data := []database.LibraryPaperPermission{}
	database.DB.Where("user_id = ?", id).
		Preload("Paper.Library").Find(&data)

	for _, d := range data {
		respBody = append(respBody, profilePermissionResponse{
			CreatedAt:    d.CreatedAt,
			Id:           d.ID,
			PaperUrl:     d.RedirectUrl,
			PaperTitle:   d.Paper.Title,
			PaperSubject: d.Paper.Subject,
			PaperIssn:    d.Paper.Issn,
			Library:      d.Paper.Library.Name,
			Purpose:      d.Purpose,
			Accepted:     d.Accepted,
		})
	}
	return
}
