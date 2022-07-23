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
	Permission interface{} `json:"permission,omitempty"`
	Dataset    interface{} `json:"dataset,omitempty"`
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
}

type profileResponse struct {
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
	IsWhatsapp   bool        `json:"isWhatsapp,omitempty"`
	Images       string      `json:"images,omitempty"`
}

type profilePermissionResponse struct {
	CreatedAt    time.Time `json:"createdAt"`
	Id           int       `json:"id"`
	PaperTitle   string    `json:"title"`
	PaperSubject []string  `json:"subject"`
	PaperType    string    `json:"type"`
	Library      string    `json:"libraryName"`
	Purpose      string    `json:"purpose"`
	Accepted     *bool     `json:"accepted"`
}

type profileCollectionBorrow struct {
	CreatedAt    time.Time  `json:"createdAt"`
	TakedAt      *time.Time `json:"takedAt"`
	ReturnedAt   *time.Time `json:"returnedAt"`
	Title        string     `json:"title"`
	SerialNumber string     `json:"serialNumber"`
	Slug         string     `json:"slug"`
	LibraryId    int        `json:"libraryId"`
	Library      string     `json:"libraryName"`
	Status       string     `json:"status"`
}

type profilePwdUpdateRequest struct {
	OldPassword    string `json:"oldPassword"`
	Password       string `json:"password"`
	RetypePassword string `json:"retypePassword"`
}

type profileEmailUpdateRequest struct {
	Email string `json:"newEmail"`
}

type profileUsernameUpdateRequest struct {
	Username *string `json:"newUsername"`
}

type profilePictureUpdateRequest struct {
	Picture []byte `json:"newPicture"`
}

type newBorrowRequest struct {
	Id int `json:"collectionId"`
}

type newPermissionRequest struct {
	Id      int    `json:"paperId"`
	Purpose string `json:"requestPurpose"`
}

type profileUpdateRequest struct {
	Name         string  `json:"name,omitempty"`
	Gender       *string `json:"gender,omitempty"`
	PlaceOfBirth *string `json:"placeOfBirth,omitempty"`
	DateOfBirth  *string `json:"dateOfBirth,omitempty"`
	Address      *string `json:"address,omitempty"`
	Institution  *string `json:"institution,omitempty"`
	Profession   *string `json:"profession,omitempty"`
	PhoneCode    *string `json:"phoneCode,omitempty"`
	PhoneNo      *string `json:"phoneNo,omitempty"`
	IsWhatsapp   bool    `json:"isWhatsapp,omitempty"`
}

type borrowDataset struct {
	Month     string `json:"month"`
	Count     int64  `json:"count"`
	Requested int64  `json:"requested"`
	Taked     int64  `json:"taked"`
	Finished  int64  `json:"finished"`
	Canceled  int64  `json:"canceled"`
}

type bookDataset struct {
	Count int64 `json:"count"`
	New   int64 `json:"new"`
	Great int64 `json:"great"`
	Good  int64 `json:"good"`
	Bad   int64 `json:"bad"`
}

type paperDataset struct {
	Count   int64 `json:"count"`
	Journal int64 `json:"journal"`
	Thesis  int64 `json:"thesis"`
	Other   int64 `json:"other"`
}

type accessDataset struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type libraryDashboardResponse struct {
	Borrow     []borrowDataset `json:"borrow"`
	Access     []accessDataset `json:"access"`
	BookStatus bookDataset     `json:"bookStatus"`
	PaperType  paperDataset    `json:"paperType"`
	Monthly    monthCount      `json:"monthly"`
}

type monthCount struct {
	Borrow int64 `json:"borrow"`
	Access int64 `json:"access"`
}

type datarange struct {
	Id       int
	FromDate string
	ToDate   string
}

type collectionUpdateRequest struct {
	Status       int `json:"status"`
	Availability int `json:"availability"`
}

type collectionAddRequests struct {
	Book       *bookRequest           `json:"book,omitempty"`
	BookSlug   *string                `json:"slug,omitempty"`
	Collection []collectionAddRequest `json:"collection"`
}

type collectionAddRequest struct {
	SerialNumber string `json:"sn"`
	Availability int    `json:"availability"`
}

type bookRequest struct {
	Title       string `json:"title"`
	Image       []byte `json:"image"`
	Author      string `json:"author"`
	ReleaseDate string `json:"releaseDate"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Country     string `json:"country"`
	Publisher   string `json:"publisher"`
	PageCount   int    `json:"pageCount"`
	Category    string `json:"category"`
}
