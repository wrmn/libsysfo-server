package database

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LibraryData struct {
	ID            int             `json:"id" gorm:"type:int;size:32;primaryKey"`
	UserID        int             `json:"userId" gorm:"type:int;size:32;autoIncrement:false"`
	User          ProfileAccount  `json:"user" gorm:"foreignKey:UserID"`
	Name          string          `json:"name" gorm:"type:varchar(64);not null"`
	Address       string          `json:"address" gorm:"type:varchar(128);not null"`
	Coordinate    pq.Float64Array `json:"coordinate" gorm:"type:decimal[]"`
	Description   string          `json:"description" gorm:"type:text"`
	ImagesMain    string          `json:"imagesMain" gorm:"type:varchar(256)"`
	ImagesContent pq.StringArray  `json:"imagesContent" gorm:"type:varchar(256)[]"`
	Webpage       string          `json:"webpage" gorm:"type:varchar(32)"`
}

// TODO: Add serial number here
type LibraryCollection struct {
	gorm.Model
	ID           int                       `json:"id" gorm:"type:int;primaryKey;size:32"`
	SerialNumber string                    `json:"serialNumber" gorm:"type:varchar(32);unique;not null"`
	LibraryID    int                       `gorm:"type:int;size:32;autoIncrement:false"`
	Library      LibraryData               `json:"library" gorm:"foreignKey:LibraryID"`
	BookID       int                       `gorm:"type:int;size:32;autoIncrement:false"`
	Book         Book                      `json:"book" gorm:"foreignKey:BookID"`
	Availability bool                      `json:"availability" gorm:"not null"`
	Status       int                       `json:"status" gorm:"type:int;size:32;not null"`
	Borrow       []LibraryCollectionBorrow `json:"borrowHistory" gorm:"foreignKey:CollectionID"`
}

type LibraryCollectionBorrow struct {
	ID           int               `json:"id" gorm:"type:int;primaryKey;size:32"`
	CreatedAt    datatypes.Date    `json:"createdAt" gorm:"type:timestamp"`
	ReturnedAt   datatypes.Date    `json:"returnedAt" gorm:"type:timestamp"`
	CollectionID int               `gorm:"type:int;size:32;autoIncrement:false"`
	Collection   LibraryCollection `json:"collection" gorm:"foreignKey:CollectionID"`
	UserID       int               `gorm:"type:int;size:32;autoIncrement:false"`
	User         ProfileAccount    `json:"user" gorm:"foreignKey:UserID"`
}

type LibraryPaper struct {
	gorm.Model
	ID          int                      `json:"id" gorm:"type:int;primaryKey;size:32"`
	LibraryID   int                      `gorm:"type:int;size:32;autoIncrement:false"`
	Library     LibraryData              `json:"library" gorm:"foreignKey:LibraryID"`
	Title       string                   `json:"title" gorm:"type:varchar(128)"`
	Subject     pq.StringArray           `json:"subject" gorm:"type:varchar(16)[]"`
	Abstract    string                   `json:"abstract"`
	Issn        string                   `json:"issn" gorm:"type:varchar(16)"`
	Description datatypes.JSON           `json:"description"`
	Access      bool                     `json:"access"`
	PaperUrl    string                   `gorm:"type:varchar(128)"`
	Permission  []LibraryPaperPermission `json:"permission" gorm:"foreignKey:PaperID"`
}

type LibraryPaperPermission struct {
	gorm.Model
	ID          int            `gorm:"type:int;primaryKey;size:32"`
	PaperID     int            `gorm:"type:int;size:32;autoIncrement:false"`
	Paper       LibraryPaper   `json:"paper" gorm:"foreignKey:PaperID"`
	UserID      int            `gorm:"type:int;size:32;autoIncrement:false"`
	User        ProfileAccount `json:"user" gorm:"foreignKey:UserID"`
	RedirectUrl string         `json:"redirectUrl" gorm:"type:varchar(128)"`
}

type LibraryPaperAccess struct {
	ID           int                    `json:"id" gorm:"type:int;primaryKey;size:32"`
	CreatedAt    datatypes.Date         `json:"createdAt" gorm:"type:timestamp"`
	PermissionID int                    `gorm:"type:int;size:32;autoIncrement:false"`
	Permission   LibraryPaperPermission `json:"paper" gorm:"foreignKey:PermissionID"`
}
