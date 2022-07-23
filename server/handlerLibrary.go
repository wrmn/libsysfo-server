package server

import (
	"libsysfo-server/database"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func allLibraries(w http.ResponseWriter, _ *http.Request) {
	result := []database.LibraryData{}
	database.DB.Find(&result)
	var libRespBody []libraryResponse

	for _, e := range result {
		libRespBody = append(libRespBody, libraryResponse{
			Id:            e.ID,
			Name:          e.Name,
			Address:       e.Address,
			Coordinate:    e.Coordinate,
			Description:   e.Description,
			ImagesMain:    e.ImagesMain,
			ImagesContent: e.ImagesContent,
		})
	}
	response{
		Data: responseBody{
			Library: libRespBody,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}

func singleLibrary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reqId, err := strconv.Atoi(vars["id"])

	if err != nil {
		badRequest(w, "invalid id")
		return
	}

	result := database.LibraryData{
		ID: reqId,
	}

	database.DB.First(&result)

	book, paper := []database.LibraryCollection{}, []database.LibraryPaper{}
	bookCount, paperCount :=
		database.DB.Where("library_id=?", result.ID).Find(&book).RowsAffected,
		database.DB.Where("library_id=?", result.ID).Find(&paper).RowsAffected

	bookData, err := getLibraryBook(reqId)
	if err != nil {
		intServerError(w, err)
		return
	}

	paperData, err := getLibraryPaper(reqId)
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Data: responseBody{
			Library: libraryResponse{
				Id:                   result.ID,
				Name:                 result.Name,
				Address:              result.Address,
				Coordinate:           result.Coordinate,
				Description:          result.Description,
				ImagesMain:           result.ImagesMain,
				ImagesContent:        result.ImagesContent,
				TotalBookCollection:  bookCount,
				TotalPaperCollection: paperCount,
			},
			Book:  &bookData,
			Paper: &paperData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}

func allPapers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))

	if page == 0 {
		page = 1
	}

	data := []database.LibraryPaper{}
	database.DB.Scopes(paginator(r, 24), paperFilter(r)).Find(&data)
	paperRespBody := []paperResponse{}
	for _, e := range data {
		paperRespBody = append(paperRespBody, paperResponse{
			Id:          e.ID,
			Title:       e.Title,
			Subject:     e.Subject,
			Type:        e.Type,
			Description: e.Description,
			Access:      e.Access,
		})
	}

	paginateData := paginate{Data: len(paperRespBody)}.generate(r, page)
	count := database.DB.Find(&data).RowsAffected
	if page*24 >= int(count) {
		paginateData.Next = ""
	}

	response{
		Data: responseBody{
			Paper:    paperRespBody,
			Paginate: &paginateData,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)

}

func singlePaper(w http.ResponseWriter, r *http.Request) {
	paperId := mux.Vars(r)["id"]

	result := database.LibraryPaper{}
	err := database.DB.Where("Id = ?", paperId).Find(&result).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	libraryResult := database.LibraryData{}
	err = database.DB.Where("Id = ?", result.LibraryID).Find(&libraryResult).Error
	if err != nil {
		intServerError(w, err)
		return
	}

	response{
		Data: responseBody{
			Paper: paperResponse{
				Id:          result.ID,
				Title:       result.Title,
				Subject:     result.Subject,
				Abstract:    result.Abstract,
				Type:        result.Type,
				Description: result.Description,
				Access:      result.Access,
			},
			Library: libraryResponse{
				Id:         libraryResult.ID,
				Name:       libraryResult.Name,
				Address:    libraryResult.Address,
				Coordinate: libraryResult.Coordinate,
			},
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Success",
	}.responseFormatter(w)
}
