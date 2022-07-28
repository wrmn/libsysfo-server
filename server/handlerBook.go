package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"libsysfo-server/database"
	bookserver "libsysfo-server/utility/book-server"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func serverEndpoint(path string) (template bookserver.BookResponse, err error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	link := fmt.Sprintf("%s%s", os.Getenv("BOOK_SERVER_URL"), path)
	token := fmt.Sprintf("Bearer %s", os.Getenv("BOOK_SERVER_TOKEN"))
	req, err := http.NewRequest("GET", link, nil)
	req.Header.Set("user-agent", "golang application")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", token)
	response, err := client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		err = errors.New("server unavailable")
		return
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &template)
	if err != nil {
		return
	}
	return
}

func allBooks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	bookRespBody := []bookResponse{}
	page, _ := strconv.Atoi(q.Get("page"))
	if page == 0 {
		page = 1
	}
	data := []database.Book{}
	if q.Has("source") {
		if q.Get("source") == "gramedia" {
			path := fmt.Sprintf("/api/books?page=%d", page)
			if q.Has("keyword") {
				path = fmt.Sprintf("/api/books/search?page=%d&keyword=%s", page, url.QueryEscape(q.Get("keyword")))
			}
			data, err := serverEndpoint(path)
			if err != nil {
				intServerError(w, err)
				return
			}

			result := data.Books
			for _, r := range result {
				bookRespBody = append(bookRespBody, bookResponse{
					Title:  *r.Title,
					Image:  *r.Image,
					Author: *r.Author,
					Slug:   *r.Slug,
					Source: "gramedia",
				})
			}
		} else if q.Get("source") == "local" {
			database.DB.Table("books").
				Scopes(bookFilter(r, 10), bookDetailFilter(r)).
				Joins("left join book_details on books.id = book_details.id").
				Scan(&bookRespBody)
		}
	} else {
		if page == 0 {
			page = 1
		}

		count := int(database.DB.
			Find(&[]database.Book{}).RowsAffected / 24)
		database.DB.Scopes(paginator(r, 24)).Preload("BookDetail").Find(&data)

		bookRespBody = localBookGenerate(data, bookRespBody)

		if page > count {
			path := fmt.Sprintf("/api/books?page=%d", page-count)
			data, err := serverEndpoint(path)
			if err != nil {
				intServerError(w, err)
				return
			}

			result := data.Books
			for _, r := range result {
				bookRespBody = append(bookRespBody, bookResponse{
					Title:  *r.Title,
					Image:  *r.Image,
					Author: *r.Author,
					Slug:   *r.Slug,
					Source: "gramedia",
				})
			}
		}
	}

	paginateData := paginate{Data: len(bookRespBody)}.generate(r, page)

	response{
		Data: responseBody{
			Book:     bookRespBody,
			Paginate: &paginateData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func singleBook(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	result := database.Book{}
	query := database.DB.Preload("BookDetail").
		Where("slug = ?", slug).Find(&result)
	exist, err := database.CheckExist(query)
	if err != nil {
		intServerError(w, err)
		return
	}
	bookRespBody := bookResponse{}
	subResult := []database.LibraryCollection{}
	libRespBody := []libraryCollectionResponse{}

	if exist != 0 {
		bookRespBody = setBookResponse(result)

		err = database.DB.
			Preload("Library").
			Where("book_id = ?", result.ID).
			Find(&subResult).Error
		if err != nil {
			intServerError(w, err)
			return
		}
		for _, e := range subResult {
			libRespBody = append(libRespBody, libraryCollectionResponse{
				Id:           e.ID,
				LibraryId:    e.LibraryID,
				Name:         e.Library.Name,
				Coordinate:   e.Library.Coordinate,
				SerialNumber: e.SerialNumber,
				Availability: e.Availability,
				Status:       e.Status,
			})
		}

	} else {
		path := fmt.Sprintf("/api/books/%s/detail", slug)
		data, err := serverEndpoint(path)
		if err != nil {
			intServerError(w, err)
			return
		}
		result := data.Book
		details := result.Detail
		bookRespBody.Title = *result.Title
		bookRespBody.Image = *result.Image
		bookRespBody.Author = *result.Author
		bookRespBody.Slug = *result.Slug
		bookRespBody.Source = "gramedia"
		bookRespBody.ReleaseDate = *details.ReleaseDate
		bookRespBody.Description = *details.Description
		bookRespBody.Language = *details.Language
		bookRespBody.Country = *details.Country
		bookRespBody.PageCount = int(*details.PageCount)
		bookRespBody.Publisher = *details.Publisher
		bookRespBody.Category = *details.Category
		bookRespBody.Origin = *result.OriginalURL
	}
	response{
		Data: responseBody{
			Book:       bookRespBody,
			Collection: libRespBody,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func localBookGenerate(data []database.Book, bookRespBody []bookResponse) []bookResponse {
	for _, d := range data {
		bookRespBody = append(bookRespBody, bookResponse{
			Title:       d.Title,
			Image:       d.Image,
			Author:      d.Author,
			Slug:        d.Slug,
			ReleaseDate: d.BookDetail.ReleaseDate,
			Description: d.BookDetail.Description,
			Language:    d.BookDetail.Language,
			Country:     d.BookDetail.Country,
			PageCount:   d.BookDetail.PageCount,
			Publisher:   d.BookDetail.Publisher,
			Category:    d.BookDetail.Category,
		})
	}

	return bookRespBody
}
