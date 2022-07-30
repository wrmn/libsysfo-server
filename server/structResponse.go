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
	Profile               interface{} `json:"profile,omitempty"`
	Library               interface{} `json:"library,omitempty"`
	Book                  interface{} `json:"book,omitempty"`
	Collection            interface{} `json:"collection,omitempty"`
	AlternativeCollection interface{} `json:"alternativeCollection,omitempty"`
	Borrow                interface{} `json:"borrow,omitempty"`
	Paper                 interface{} `json:"paper,omitempty"`
	Permission            interface{} `json:"permission,omitempty"`
	Access                interface{} `json:"Access,omitempty"`
	Dataset               interface{} `json:"dataset,omitempty"`
	User                  interface{} `json:"user,omitempty"`
	Token                 string      `json:"token,omitempty"`
	Paginate              *paginate   `json:"paginate,omitempty"`
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
	Id          int                        `json:"id,omitempty"`
	Title       string                     `json:"title"`
	Image       string                     `json:"image"`
	Author      string                     `json:"author"`
	Slug        string                     `json:"slug"`
	ReleaseDate string                     `json:"releaseDate,omitempty"`
	Description string                     `json:"description,omitempty"`
	Language    string                     `json:"language,omitempty"`
	Country     string                     `json:"country,omitempty"`
	Publisher   string                     `json:"publisher,omitempty"`
	PageCount   int                        `json:"pageCount,omitempty"`
	Category    string                     `json:"category,omitempty"`
	Origin      string                     `json:"origin,omitempty"`
	Source      string                     `json:"source,omitempty"`
	Status      *libraryCollectionResponse `json:"status,omitempty" gorm:"type:text"`
}

type libraryCollectionResponse struct {
	Id           int       `json:"id"`
	SerialNumber string    `json:"sn"`
	Name         string    `json:"name,omitempty"`
	LibraryId    int       `json:"libraryId,omitempty"`
	Coordinate   []float64 `json:"coordinate,omitempty"`
	Availability int       `json:"availability"`
	Status       int       `json:"status"`
}

type profileCollectionBorrowResponse struct {
	BorrowId     int        `json:"id,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	AcceptedAt   *time.Time `json:"acceptedAt"`
	TakedAt      *time.Time `json:"takedAt"`
	ReturnedAt   *time.Time `json:"returnedAt"`
	CanceledAt   *time.Time `json:"canceledAt"`
	Title        string     `json:"title"`
	SerialNumber string     `json:"serialNumber"`
	CollectionId int        `json:"collectionId"`
	Slug         string     `json:"slug"`
	LibraryId    int        `json:"libraryId,omitempty"`
	Library      string     `json:"libraryName,omitempty"`
	UserId       int        `json:"userId,omitempty"`
	UserName     string     `json:"userName,omitempty"`
	Status       string     `json:"status"`
}

type adminInformationResponse struct {
	Username      string          `json:"username"`
	Email         string          `json:"email"`
	Library       string          `json:"libraryName"`
	Image         string          `json:"libraryImage"`
	Address       string          `json:"libraryAddress"`
	Coordinate    pq.Float64Array `json:"coordinate"`
	Description   string          `json:"description"`
	ContentImages pq.StringArray  `json:"contentImages"`
	Webpage       string          `json:"webpage"`
}

type paperResponse struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Subject     pq.StringArray `json:"subject"`
	Abstract    string         `json:"abstract,omitempty"`
	Type        string         `json:"type"`
	Description datatypes.JSON `json:"description"`
	Access      bool           `json:"access"`
	PaperUrl    *string        `json:"paperUrl"`
}

type profileResponse struct {
	Id           *int        `json:"id,omitempty"`
	Username     *string     `json:"username,omitempty"`
	Email        string      `json:"email,omitempty"`
	Verified     interface{} `json:"verivied,omitempty"`
	Name         string      `json:"name,omitempty"`
	Gender       *string     `json:"gender,omitempty"`
	PlaceOfBirth *string     `json:"placeOfBirth,omitempty"`
	DateOfBirth  *string     `json:"dateOfBirth,omitempty"`
	Address      *string     `json:"address,omitempty"`
	Institution  *string     `json:"institution,omitempty"`
	Profession   *string     `json:"profession,omitempty"`
	PhoneCode    *string     `json:"phoneCode,omitempty"`
	PhoneNo      *string     `json:"phoneNo,omitempty"`
	IsWhatsapp   *bool       `json:"isWhatsapp,omitempty"`
	Images       string      `json:"images,omitempty"`
}

type profilePermissionResponse struct {
	CreatedAt    time.Time  `json:"createdAt"`
	AcceptedAt   *time.Time `json:"acceptedAt"`
	Id           int        `json:"id"`
	PaperTitle   string     `json:"title"`
	PaperSubject []string   `json:"subject"`
	PaperType    string     `json:"type"`
	Library      string     `json:"libraryName"`
	Purpose      string     `json:"purpose"`
}

type accessResponse struct {
	Total     int         `json:"total"`
	CreatedAt []time.Time `json:"time"`
}

type libraryDashboardResponse struct {
	Borrow     []borrowDataset `json:"borrow"`
	Access     []accessDataset `json:"access"`
	BookStatus bookDataset     `json:"bookStatus"`
	PaperType  paperDataset    `json:"paperType"`
	Monthly    monthCount      `json:"monthly"`
}
