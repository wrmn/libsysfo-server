package database

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LibraryData struct {
	UserId      int             `json:"userId" gorm:"type:int;size:32;primaryKey;autoIncrement:false"`
	User        ProfileAccount  `json:"user" gorm:"foreignKey:UserId"`
	Name        string          `json:"name" gorm:"type:varchar(32);not null"`
	Address     string          `json:"address" gorm:"type:varchar(64);not null"`
	Coordinate  pq.Float64Array `json:"coordinate" gorm:"type:decimal[]"`
	Description string          `json:"description" gorm:"type:text"`
	Images      datatypes.JSON  `json:"images" gorm:"type:json"`
	Webpage     string          `json:"webpage" gorm:"type:varchar(32)"`
}

type LibraryCollection struct {
	gorm.Model
	Id           int                       `json:"id" gorm:"type:int;primaryKey;size:32"`
	LibraryId    int                       `gorm:"type:int;size:32;autoIncrement:false"`
	Library      LibraryData               `json:"library" gorm:"foreignKey:LibraryId"`
	BookId       int                       `gorm:"type:int;size:32;autoIncrement:false"`
	Book         Book                      `json:"book" gorm:"foreignKey:BookId"`
	Availability bool                      `json:"availability" gorm:"not null"`
	Status       int                       `json:"status" gorm:"type:int;size:32;not null"`
	Borrow       []LibraryCollectionBorrow `json:"borrowHistory" gorm:"foreignKey:CollectionId"`
}

type LibraryCollectionBorrow struct {
	Id           int               `json:"id" gorm:"type:int;primaryKey;size:32"`
	CreatedAt    datatypes.Time    `json:"createdAt" gorm:"type:timestamp"`
	CollectionId int               `gorm:"type:int;size:32;autoIncrement:false"`
	Collection   LibraryCollection `json:"collection" gorm:"foreignKey:CollectionId"`
	UserId       int               `gorm:"type:int;size:32;autoIncrement:false"`
	User         ProfileAccount    `json:"user" gorm:"foreignKey:UserId"`
}

type LibraryPaper struct {
	gorm.Model
	Id          int                      `json:"id" gorm:"type:int;primaryKey;size:32"`
	LibraryId   int                      `gorm:"type:int;size:32;autoIncrement:false"`
	Library     LibraryData              `json:"library" gorm:"foreignKey:LibraryId"`
	Title       string                   `json:"title" gorm:"type:varchar(32)"`
	Subject     pq.StringArray           `json:"subject" gorm:"type:varchar(16)[]"`
	Abstract    string                   `json:"abstract"`
	Issn        string                   `json:"issn" gorm:"type:varchar(16)"`
	Description datatypes.JSON           `json:"description"`
	Access      bool                     `json:"access"`
	PaperUrl    string                   `gorm:"type:varchar(128)"`
	Permission  []LibraryPaperPermission `json:"permission" gorm:"foreignKey:PaperId"`
}

type LibraryPaperPermission struct {
	gorm.Model
	Id          int            `gorm:"type:int;primaryKey;size:32"`
	PaperId     int            `gorm:"type:int;size:32;autoIncrement:false"`
	Paper       LibraryPaper   `json:"paper" gorm:"foreignKey:PaperId"`
	UserId      int            `gorm:"type:int;size:32;autoIncrement:false"`
	User        ProfileAccount `json:"user" gorm:"foreignKey:UserId"`
	RedirectUrl string         `json:"redirectUrl" gorm:"type:varchar(32)"`
}

type LibraryPaperAccess struct {
	Id           int                    `json:"id" gorm:"type:int;primaryKey;size:32"`
	CreatedAt    datatypes.Time         `json:"createdAt" gorm:"type:timestamp"`
	PermissionId int                    `gorm:"type:int;size:32;autoIncrement:false"`
	Permission   LibraryPaperPermission `json:"paper" gorm:"foreignKey:PermissionId"`
}
