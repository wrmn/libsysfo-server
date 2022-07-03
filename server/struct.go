package server

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type response struct {
	Data        responseBody `json:"data,omitempty"`
	Status      int          `json:"status"`
	Reason      string       `json:"reason"`
	Description string       `json:"description"`
}

type responseBody struct {
	Profile    interface{} `json:"profile,omitempty"`
	Library    interface{} `json:"library,omitempty"`
	Book       interface{} `json:"book,omitempty"`
	Collection interface{} `json:"collection,omitempty"`
	Paper      interface{} `json:"paper,omitempty"`
	Token      string      `json:"token,omitempty"`
	Paginate   *paginate   `json:"paginate,omitempty"`
}

type paginate struct {
	Data    int    `json:"dataTotal"`
	Next    string `json:"nextPage,omitempty"`
	Current string `json:"currentPage"`
	Prev    string `json:"previousPage"`
}

type libraryResponse struct {
	Id                   int              `json:"id"`
	Name                 string           `json:"name"`
	Address              string           `json:"address"`
	Coordinate           []float64        `json:"coordinate"`
	Description          string           `json:"description"`
	ImagesMain           string           `json:"imagesMain"`
	ImagesContent        []string         `json:"imagesContent"`
	TotalBookCollection  int64            `json:"totalBookCollection,omitempty"`
	BookCollection       *[]bookResponse  `json:"bookCollection,omitempty"`
	TotalPaperCollection int64            `json:"totalPaperCollection,omitempty"`
	PaperCollection      *[]paperResponse `json:"paperCollection,omitempty"`
}

type bookResponse struct {
	Title       string                       `json:"title"`
	Image       string                       `json:"image"`
	Author      string                       `json:"author"`
	Slug        string                       `json:"slug"`
	ReleaseDate string                       `json:"releaseDate,omitempty"`
	Description string                       `json:"description,omitempty"`
	Language    string                       `json:"language,omitempty"`
	Country     string                       `json:"country,omitempty"`
	Publisher   string                       `json:"publisher,omitempty"`
	PageCount   float64                      `json:"pageCount,omitempty"`
	Category    string                       `json:"category,omitempty"`
	Origin      string                       `json:"origin,omitempty"`
	AvailableOn *[]libraryCollectionResponse `json:"availableOn,omitempty"`
	Status      *libraryCollectionResponse   `json:"status,omitempty"`
}

type libraryCollectionResponse struct {
	SerialNumber string    `json:"sn"`
	Name         string    `json:"name,omitempty"`
	LibraryId    int       `json:"libraryId,omitempty"`
	Coordinate   []float64 `json:"coordinate,omitempty"`
	Availability bool      `json:"availability"`
	Status       int       `json:"status"`
}

type paperResponse struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Subject     pq.StringArray `json:"subject"`
	Abstract    string         `json:"abstract,omitempty"`
	Issn        string         `json:"issn"`
	Description datatypes.JSON `json:"description"`
	Access      bool           `json:"access"`
}
