package server

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"libsysfo-server/database"
	"libsysfo-server/utility/cred"
	"net/http"
	"regexp"
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

func pwdLocator(r *http.Request) (pwd string, err error) {
	pwdHead := r.Header.Values("Account-auth")
	if len(pwdHead) == 0 {
		err = errors.New("authorization required")
		return
	}
	pwd = pwdHead[0]
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

func findUser(cred *cred.TokenModel, pwd string) (user database.ProfileAccount, err error) {
	password := fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	result := database.DB.
		Where("email = ? AND password = ?", cred.Email, password).
		Or("username = ? AND password = ?", cred.Username, password).
		Find(&user)

	if result.RowsAffected == 0 {
		err = errors.New("invalid password")
		return
	}
	return
}

func getUser(cred *cred.TokenModel) (user database.ProfileAccount, err error) {
	err = database.DB.
		Where("email = ? ", cred.Email).
		Or("username = ? ", cred.Username).
		Find(&user).Error
	if user.AccountType != 3 {
		err = errors.New("user not allowed")
		return
	}
	return
}

func findUserData(id int) (user database.ProfileData, err error) {
	result := database.DB.
		Where("user_id = ?", id).
		Find(&user)

	if result.RowsAffected == 0 {
		err = errors.New("user not found")
		return
	}
	return
}

func userStats(id int) (stats bool, err error) {
	res, err := findUserData(id)
	if res.Address1 == nil ||
		res.DateOfBirth == nil ||
		res.Gender == nil ||
		res.Institution == nil ||
		res.PhoneCode == nil ||
		res.PlaceOfBirth == nil ||
		res.Profession == nil ||
		res.VerifiedAt == nil {
		stats = false
		return
	}
	return true, err
}

func handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response{
			Status:      http.StatusNotFound,
			Reason:      "Not Found",
			Description: "Request not found",
		}.responseFormatter(w)
	})
}

func handleNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response{
			Status:      http.StatusMethodNotAllowed,
			Reason:      "Method Not Allowed",
			Description: "Request not allowed with this method",
		}.responseFormatter(w)
	})
}

func getLoginData(r *http.Request) (user database.ProfileAccount, err error) {
	var e cred.FormAuth
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()
	err = decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			err = errors.New("Wrong Type provided for field " + unmarshalErr.Field)
		}
		return
	}

	e.Password = fmt.Sprintf("%x", md5.Sum([]byte(e.Password)))

	user = database.ProfileAccount{}
	result := database.DB.Where("email = ? AND password = ?", e.Indicator, e.Password).Or("username = ? AND password = ?", e.Indicator, e.Password).Find(&user).RowsAffected
	if result == 0 {
		err = errors.New("invalid username or password ")
	}

	return
}

func getLibraryBook(libId int) (bookData []bookResponse, err error) {
	bookQuery := []database.LibraryCollection{}
	err = database.DB.Where("library_id = ?", libId).Preload("Book", func(db *gorm.DB) *gorm.DB {
		return database.DB.Preload("BookDetail")
	}).Order("id desc").Find(&bookQuery).Error

	if err != nil {
		return
	}

	for _, d := range bookQuery {
		bookData = append(bookData, bookResponse{
			Id:          d.ID,
			Title:       d.Book.Title,
			Image:       d.Book.Image,
			Author:      d.Book.Author,
			Slug:        d.Book.Slug,
			ReleaseDate: d.Book.BookDetail.ReleaseDate,
			Language:    d.Book.BookDetail.Language,
			Description: d.Book.BookDetail.Description,
			Country:     d.Book.BookDetail.Country,
			Publisher:   d.Book.BookDetail.Publisher,
			PageCount:   d.Book.BookDetail.PageCount,
			Category:    d.Book.BookDetail.Category,
			Status: &libraryCollectionResponse{
				Availability: d.Availability,
				Status:       d.Status,
			},
		})
	}

	return
}

func getLibraryPaper(libId int) (paperData []paperResponse, err error) {
	paperQuery := []database.LibraryPaper{}
	err = database.DB.Where("library_id = ?", libId).Find(&paperQuery).Error
	if err != nil {
		return
	}

	for _, c := range paperQuery {
		paperData = append(paperData, paperResponse{
			Id:          c.ID,
			Title:       c.Title,
			Subject:     c.Subject,
			Abstract:    c.Abstract,
			Type:        c.Type,
			Description: c.Description,
			Access:      c.Access,
		})
	}
	return
}

func setBookResponse(result database.Book) (bookRespBody bookResponse) {
	bookRespBody.Title = result.Title
	bookRespBody.Image = result.Image
	bookRespBody.Author = result.Author
	bookRespBody.Slug = result.Slug
	bookRespBody.Source = result.Source
	bookRespBody.ReleaseDate = result.BookDetail.ReleaseDate
	bookRespBody.Description = result.BookDetail.Description
	bookRespBody.Language = result.BookDetail.Language
	bookRespBody.Country = result.BookDetail.Country
	bookRespBody.PageCount = int(result.BookDetail.PageCount)
	bookRespBody.Publisher = result.BookDetail.Publisher
	bookRespBody.Category = result.BookDetail.Category
	return
}

func checkExist(q *gorm.DB) (int64, error) {
	return q.RowsAffected, q.Error
}

func slugGenerator(title string) (slug string) {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	slug = re.ReplaceAllString(strings.ToLower(title), "-")
	var exist int64 = 1
	rep := 0
	for exist != 0 {
		query := database.DB.Preload("BookDetail").
			Where("slug = ?", slug).Find(&database.Book{})
		exist, _ = checkExist(query)
		if exist != 0 {
			rep += 1
			slug = fmt.Sprintf("%s-%d", slug, rep)
		}
	}

	return
}
