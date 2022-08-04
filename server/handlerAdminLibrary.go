package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"libsysfo-server/database"
	"libsysfo-server/utility/imgkit"
	"net/http"
	"strconv"
	"time"
)

func adminLogin(w http.ResponseWriter, r *http.Request) {
	user, err := getLoginData(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	} else if user.AccountType != 2 {
		badRequest(w, "user not allowed")
		return
	}
	loginHandler(w, user)
}

func adminInformation(w http.ResponseWriter, r *http.Request) {
	data, invalid := checkToken(r, w)
	if invalid {
		return
	}

	if data.AccountType != 2 {
		unauthorizedRequest(w, errors.New("user not allowed"))
		return
	}

	libraryData, invalid := getLibraryData(data.ID, w)
	if invalid {
		return
	}

	response{
		Data: responseBody{
			Profile: adminInformationResponse{
				Username:      *data.Username,
				Email:         data.Email,
				Library:       libraryData.Name,
				Image:         libraryData.ImagesMain,
				Address:       libraryData.Address,
				Coordinate:    libraryData.Coordinate,
				Description:   libraryData.Description,
				ContentImages: libraryData.ImagesContent,
				Webpage:       libraryData.Webpage,
			},
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func libraryDashboard(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	setBody, err := libraryDashboardResponse{}.fill(libraryData.ID, r)
	if err != nil {
		badRequest(w, "invalid date range")
	}

	response{
		Data: responseBody{
			Dataset: setBody,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func (data libraryDashboardResponse) fill(libId int, r *http.Request) (libraryDashboardResponse, error) {
	q := r.URL.Query()
	queryParam := datarange{
		Id:       libId,
		FromDate: time.Now().AddDate(0, -6, 0).Format("20060102"),
		ToDate:   time.Now().AddDate(0, 1, 0).Format("20060102"),
	}

	if q.Has("from") && q.Has("to") {
		queryParam.FromDate = q.Get("from")
		t, err := time.Parse("20060102", q.Get("to"))
		if err != nil {
			return libraryDashboardResponse{}, err
		}
		queryParam.ToDate = t.AddDate(0, 1, 0).Format("20060102")
	}

	return libraryDashboardResponse{
		Borrow:     getBorrowDataset(queryParam),
		BookStatus: getBookDataset(libId),
		Access:     getAccessDataset(queryParam),
		PaperType:  getPaperDataset(libId),
		Monthly:    fill(libId),
	}, nil
}

func fill(libId int) (data monthCount) {
	lastMonth := time.Now().AddDate(0, -1, 0).Format("20060102")
	today := time.Now().Format("20060102")
	borrowRows, _ := database.DB.Raw(`SELECT
		COUNT(*) as count
  	FROM library_collection_borrows
	LEFT JOIN library_collections 
	ON library_collections.id=library_collection_borrows.collection_id
	WHERE  library_id = ? 
		AND library_collection_borrows.created_at BETWEEN ? AND ?`, libId, lastMonth, today).Rows()
	defer borrowRows.Close()
	for borrowRows.Next() {
		database.DB.ScanRows(borrowRows, &data.Borrow)
	}

	accessRows, _ := database.DB.Raw(`SELECT
	COUNT(*) as count
	FROM library_paper_accesses
	LEFT JOIN library_paper_permissions 
		ON library_paper_permissions.id=library_paper_accesses.permission_id
	LEFT JOIN library_papers 
		ON library_papers.id=library_paper_permissions.paper_id
	WHERE library_papers.library_id = ?
	AND library_paper_accesses.created_at BETWEEN ? AND ?`, libId, lastMonth, today).Rows()
	defer accessRows.Close()
	for accessRows.Next() {
		database.DB.ScanRows(accessRows, &data.Access)
	}
	return
}

func getPaperDataset(libId int) (paperBody paperDataset) {
	paperRows, _ := database.DB.Raw(`SELECT
		COUNT(*) as count,
		SUM(CASE WHEN type = 'journal' THEN 1 ELSE 0 END) as journal,
		SUM(CASE WHEN type = 'thesis' THEN 1 ELSE 0 END) as thesis,
		SUM(CASE WHEN type = 'other document' THEN 1 ELSE 0 END) as other
  	FROM library_papers
	WHERE library_id = ?`, libId).Rows()

	defer paperRows.Close()
	for paperRows.Next() {
		database.DB.ScanRows(paperRows, &paperBody)
	}
	return paperBody
}

func getBookDataset(libId int) (bookBody bookDataset) {
	bookRows, _ := database.DB.Raw(`SELECT
		COUNT(*) as count,
		SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) as new,
		SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) as great,
		SUM(CASE WHEN status = 3 THEN 1 ELSE 0 END) as good,
		SUM(CASE WHEN status = 4 THEN 1 ELSE 0 END) as bad
	FROM library_collections
	WHERE library_id = ?`, libId).Rows()

	defer bookRows.Close()
	for bookRows.Next() {
		database.DB.ScanRows(bookRows, &bookBody)
	}
	return bookBody
}

func getAccessDataset(q datarange) (accessBody []accessDataset) {
	accessResult := accessDataset{}
	accessRows, _ := database.DB.Raw(`SELECT
		COUNT(*) as count,
		date_trunc('month', library_paper_accesses.created_at) as month
	FROM library_paper_accesses
	LEFT JOIN library_paper_permissions 
		ON library_paper_permissions.id=library_paper_accesses.permission_id
	LEFT JOIN library_papers 
		ON library_papers.id=library_paper_permissions.paper_id
	WHERE library_papers.library_id = ?
		AND library_paper_accesses.created_at BETWEEN ? AND ?
	GROUP BY month 
	ORDER BY month`, q.Id, q.FromDate, q.ToDate).Rows()
	defer accessRows.Close()
	for accessRows.Next() {
		database.DB.ScanRows(accessRows, &accessResult)
		accessBody = append(accessBody, accessResult)
	}
	return accessBody
}

func getBorrowDataset(q datarange) (borrowBody []borrowDataset) {
	borrowResult := borrowDataset{}
	borrowRows, _ := database.DB.Raw(`SELECT
		COUNT(*) as count,
		date_trunc('month', library_collection_borrows.created_at) as month,
		SUM(CASE WHEN 
			library_collection_borrows.returned_at IS NOT NULL 
		THEN 1 ELSE 0 END) as finished,
		SUM(CASE WHEN 
			library_collection_borrows.canceled_at IS NOT NULL
		THEN 1 ELSE 0 END) as canceled
  	FROM library_collection_borrows
	LEFT JOIN library_collections 
	ON library_collections.id=library_collection_borrows.collection_id
	WHERE library_collections.library_id=?
		AND library_collection_borrows.created_at BETWEEN ? AND ?
  	GROUP BY month
	ORDER BY month`, q.Id, q.FromDate, q.ToDate).Rows()
	defer borrowRows.Close()
	for borrowRows.Next() {
		database.DB.ScanRows(borrowRows, &borrowResult)
		borrowBody = append(borrowBody, borrowResult)
	}
	return borrowBody
}

func libraryImage(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	var e libraryImageUpdateRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	img := imgkit.ImgInformation{
		File:     e.File,
		FileName: strconv.Itoa(libraryData.ID),
		Folder:   fmt.Sprintf("/book/%d/", libraryData.ID),
	}

	upr, err := img.UploadImage()
	if err != nil {
		intServerError(w, err)
		return
	}

	libraryData.ImagesMain = upr.URL
	err = database.DB.Save(&libraryData).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Image updated",
	}.responseFormatter(w)
}

func libraryGeneral(w http.ResponseWriter, r *http.Request) {
	libraryData, invalid := isLibraryAdmin(w, r)
	if invalid {
		return
	}

	var e libraryGeneralUpdateRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}

	libraryData.Name = e.Name
	libraryData.Address = e.Address
	libraryData.Webpage = e.Webpage
	libraryData.Description = e.Description
	err = database.DB.Save(&libraryData).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Information updated",
	}.responseFormatter(w)
}
