package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"libsysfo-server/utility/cred"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func (data response) responseFormatter(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(data.Status)
	jsonResp, _ := json.Marshal(data)
	w.Write(jsonResp)
}

func unauthorizedRequest(w http.ResponseWriter, err error) {
	response{
		Status:      http.StatusUnauthorized,
		Reason:      "Unauthorized",
		Description: err.Error(),
	}.responseFormatter(w)
}

func badRequest(w http.ResponseWriter, msg string) {
	response{
		Status:      http.StatusBadRequest,
		Reason:      "Bad Request",
		Description: msg,
	}.responseFormatter(w)
}

func intServerError(w http.ResponseWriter, err error) {
	response{
		Status:      http.StatusInternalServerError,
		Reason:      "Internal Server Error",
		Description: err.Error(),
	}.responseFormatter(w)
}

func authVerification(r *http.Request) (tokenData *jwt.Token, err error) {
	tokenHeader := r.Header.Values("Authorization")
	if len(tokenHeader) == 0 {
		err = errors.New("authorization required")
		return
	}
	token := strings.Split(tokenHeader[0], " ")
	if token[0] != "Bearer" {
		err = errors.New("need bearer authorization")
		return
	}
	tokenData, err = cred.VerifyToken(token[1])
	return

}

func paginator(r *http.Request, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func bookFilter(r *http.Request, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
		if q.Has("keyword") {
			keyword := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("keyword")))
			db = db.Where("LOWER(title) like ?", keyword).
				Or("LOWER(author) like ?", keyword)
		}
		return db
	}
}

func bookDetailFilter(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()

		if q.Has("language") {
			if q.Has("category") {
				category := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("category")))
				db = db.Where("LOWER(language) = ? AND LOWER(category) = ?", strings.ToLower(q.Get("language")), category)
				return db
			}
			db = db.Where("LOWER(language) = ? ", strings.ToLower(q.Get("language")))
			return db
		} else if q.Has("category") {
			category := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("category")))
			db = db.Where("LOWER(category) like ? ", category)
			return db
		}
		return db
	}
}

func paperFilter(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		if q.Has("keyword") {
			keyword := fmt.Sprintf("%%%s%%", strings.ToLower(q.Get("keyword")))
			db = db.Where("LOWER(title) like ?", keyword).
				Or("array_to_string(subject, ',', ' ') like ?", keyword).
				Or("description::TEXT like ?", keyword)
		}
		return db
	}
}

func (data paginate) generate(r *http.Request, page int) (result paginate) {
	// NOTE:change to https when deployed
	link := "http://" + r.Host + r.URL.Path + "?page="
	if page > 1 {
		result.Prev = link + strconv.Itoa(page-1)
	}
	result.Current = link + strconv.Itoa(page)
	result.Data = data.Data
	if result.Data != 0 {
		result.Next = link + strconv.Itoa(page+1)
	}
	return
}
