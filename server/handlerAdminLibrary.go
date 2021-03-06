package server

import (
	"errors"
	"libsysfo-server/database"
	"libsysfo-server/utility/cred"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func adminLogin(w http.ResponseWriter, r *http.Request) {
	user, err := getLoginData(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	} else if user.AccountType != 2 {
		err := errors.New("user not allowed")
		unauthorizedRequest(w, err)
		return
	}
	loginHandler(w, user)
}

func adminInformation(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	data, err := adminData(tokenData)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	response{
		Data: responseBody{
			Profile: adminInformationResponse{
				Username:      *data.Username,
				Email:         data.Email,
				Library:       data.Library.Name,
				Image:         data.Library.ImagesMain,
				Address:       data.Library.Address,
				Coordinate:    data.Library.Coordinate,
				Description:   data.Library.Description,
				ContentImages: data.Library.ImagesContent,
				Webpage:       data.Library.Webpage,
			},
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func libraryDashboard(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authVerification(r)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}
	data, err := adminData(tokenData)
	if err != nil {
		unauthorizedRequest(w, err)
		return
	}

	setBody, err := libraryDashboardResponse{}.fill(data.Library.ID, r)
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
		SUM(CASE WHEN library_collection_borrows.status = 'requested' THEN 1 ELSE 0 END) as requested,
		SUM(CASE WHEN library_collection_borrows.status = 'taked' THEN 1 ELSE 0 END) as taked,
		SUM(CASE WHEN library_collection_borrows.status = 'finished' THEN 1 ELSE 0 END) as finished,
		SUM(CASE WHEN library_collection_borrows.status = 'canceled' THEN 1 ELSE 0 END) as canceled
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

func adminData(tokenData *jwt.Token) (data database.ProfileAccount, err error) {
	cred := tokenData.Claims.(*cred.TokenModel)
	count := database.DB.Where("email = ?", cred.Email).Or("username = ?", cred.Username).
		Preload("Library").First(&data).RowsAffected
	if count != 1 {
		return data, errors.New("user not found")
	}
	return data, nil
}
