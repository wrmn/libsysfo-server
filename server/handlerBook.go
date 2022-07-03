package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"libsysfo-server/database"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func serverEndpoint(path string) (template interface{}, err error) {
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

			result := data.(map[string]interface{})["books"].([]interface{})
			for _, r := range result {
				d := r.(map[string]interface{})
				bookRespBody = append(bookRespBody, bookResponse{
					Title:  d["title"].(string),
					Image:  d["image"].(string),
					Author: d["author"].(string),
					Slug:   d["slug"].(string),
				})
			}
		} else if q.Get("source") == "local" {

			database.DB.Scopes(bookFilter(r, 10)).Preload("BookDetail").Find(&data)
			bookRespBody = localBookGenerate(data, bookRespBody)
		}
	} else {
		if page == 0 {
			page = 1
		}

		count := int(database.DB.
			Find(&[]database.Book{}).RowsAffected / 24)
		data := []database.Book{}
		database.DB.Scopes(paginator(r, 24)).Preload("BookDetail").Find(&data)

		bookRespBody = localBookGenerate(data, bookRespBody)

		if page > count {
			path := fmt.Sprintf("/api/books?page=%d", page-count)
			data, err := serverEndpoint(path)
			if err != nil {
				intServerError(w, err)
				return
			}

			result := data.(map[string]interface{})["books"].([]interface{})
			for _, r := range result {
				d := r.(map[string]interface{})
				bookRespBody = append(bookRespBody, bookResponse{
					Title:  d["title"].(string),
					Image:  d["image"].(string),
					Author: d["author"].(string),
					Slug:   d["slug"].(string),
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
	err := query.Error
	if err != nil {
		intServerError(w, err)
		return
	}
	exist := query.RowsAffected
	bookRespBody := bookResponse{}
	if exist != 0 {
		bookRespBody.Title = result.Title
		bookRespBody.Image = result.Image
		bookRespBody.Author = result.Author
		bookRespBody.Slug = result.Slug
		bookRespBody.ReleaseDate = result.BookDetail.ReleaseDate
		bookRespBody.Description = result.BookDetail.Description
		bookRespBody.Language = result.BookDetail.Language
		bookRespBody.Country = result.BookDetail.Country
		bookRespBody.PageCount = result.BookDetail.PageCount
		bookRespBody.Publisher = result.BookDetail.Publisher
		bookRespBody.Category = result.BookDetail.Category

		subResult := []database.LibraryCollection{}
		err = database.DB.
			Preload("Library").
			Where("book_id = ?", result.ID).
			Find(&subResult).Error
		if err != nil {
			intServerError(w, err)
			return
		}
		var libRespBody []libraryCollectionResponse
		for _, e := range subResult {
			libRespBody = append(libRespBody, libraryCollectionResponse{
				LibraryId:    e.LibraryID,
				Name:         e.Library.Name,
				Coordinate:   e.Library.Coordinate,
				SerialNumber: e.SerialNumber,
				Availability: e.Availability,
				Status:       e.Status,
			})
		}

		bookRespBody.AvailableOn = &libRespBody
	} else {
		path := fmt.Sprintf("/api/books/%s/detail", slug)
		data, err := serverEndpoint(path)
		if err != nil {
			intServerError(w, err)
			return
		}
		result := data.(map[string]interface{})["book"].(map[string]interface{})
		details := result["detail"].(map[string]interface{})
		bookRespBody.Title = result["title"].(string)
		bookRespBody.Image = result["image"].(string)
		bookRespBody.Author = result["author"].(string)
		bookRespBody.Slug = result["slug"].(string)
		bookRespBody.ReleaseDate = details["release_date"].(string)
		bookRespBody.Description = details["description"].(string)
		bookRespBody.Language = details["language"].(string)
		bookRespBody.Country = details["country"].(string)
		bookRespBody.PageCount = details["page_count"].(float64)
		bookRespBody.Publisher = details["publisher"].(string)
		bookRespBody.Category = details["category"].(string)
		bookRespBody.Origin = result["original_url"].(string)
	}
	response{
		Data: responseBody{
			Book: bookRespBody,
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
