package server

import (
	"time"

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
	Borrow     interface{} `json:"borrow,omitempty"`
	Permission interface{} `json:"access,omitempty"`
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
	Id                   int       `json:"id"`
	Name                 string    `json:"name"`
	Address              string    `json:"address"`
	Coordinate           []float64 `json:"coordinate"`
	Description          string    `json:"description,omitempty"`
	ImagesMain           string    `json:"imagesMain,omitempty"`
	ImagesContent        []string  `json:"imagesContent,omitempty"`
	TotalBookCollection  int64     `json:"totalBookCollection,omitempty"`
	TotalPaperCollection int64     `json:"totalPaperCollection,omitempty"`
}

type bookResponse struct {
	Title       string                     `json:"title"`
	Image       string                     `json:"image"`
	Author      string                     `json:"author"`
	Slug        string                     `json:"slug"`
	ReleaseDate string                     `json:"releaseDate,omitempty"`
	Description string                     `json:"description,omitempty"`
	Language    string                     `json:"language,omitempty"`
	Country     string                     `json:"country,omitempty"`
	Publisher   string                     `json:"publisher,omitempty"`
	PageCount   float64                    `json:"pageCount,omitempty"`
	Category    string                     `json:"category,omitempty"`
	Origin      string                     `json:"origin,omitempty"`
	Source      string                     `json:"source,omitempty"`
	Status      *libraryCollectionResponse `json:"status,omitempty" gorm:"type:text"`
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

type profileResponse struct {
	Username     string      `json:"username,omitempty"`
	Email        string      `json:"email,omitempty"`
	Verified     interface{} `json:"verivied,omitempty"`
	Name         string      `json:"name,omitempty"`
	Gender       *string     `json:"gender,omitempty"`
	PlaceOfBirth *string     `json:"placeOfBirth,omitempty"`
	DateOfBirth  *string     `json:"dateOfBirth,omitempty"`
	Address      *string     `json:"address,omitempty"`
	Institution  *string     `json:"institution,omitempty"`
	Profession   *string     `json:"profession,omitempty"`
	PhoneNo      *string     `json:"phoneNo,omitempty"`
	IsWhatsapp   bool        `json:"isWhatsapp,omitempty"`
	Images       string      `json:"images,omitempty"`
}

type profilePermissionResponse struct {
	CreatedAt    time.Time `json:"createdAt"`
	PaperUrl     string    `json:"redirectUrl,omitempty"`
	PaperTitle   string    `json:"title"`
	PaperSubject []string  `json:"subject"`
	PaperIssn    string    `json:"issn"`
	Library      string    `json:"libraryName"`
}
