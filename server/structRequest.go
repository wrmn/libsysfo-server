package server

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

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

type borrowRequest struct {
	State        string `json:"state"`
	BorrowId     *int   `json:"borrowId"`
	UserId       *int   `json:"userId"`
	CollectionId *int   `json:"collectionId"`
}

type permissionRequest struct {
	PermissionId int    `json:"permissionId"`
	State        string `json:"state"`
}

type paperAddRequest struct {
	Title       string         `json:"title"`
	Subject     pq.StringArray `json:"subject"`
	Abstract    string         `json:"abstract"`
	Type        string         `json:"type"`
	Description datatypes.JSON `json:"description"`
	Access      *bool          `json:"access"`
	PaperFile   []byte         `json:"paperFile"`
}

type fileUpdateRequest struct {
	File []byte `json:"file"`
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

type collectionUpdateRequest struct {
	SerialNumber string `json:"sn"`
	Status       int    `json:"status"`
	Availability int    `json:"availability"`
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

type libraryGeneralUpdateRequest struct {
	Name        string `json:"name"`
	Webpage     string `json:"webpage"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

type libraryImageUpdateRequest struct {
	File []byte `json:"file"`
}
