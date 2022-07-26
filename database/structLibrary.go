package database

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LibraryData struct {
	ID            int             `json:"id" gorm:"type:int;size:32;primaryKey"`
	UserID        int             `json:"userId" gorm:"type:int;size:32;autoIncrement:false"`
	Name          string          `json:"name" gorm:"type:varchar(64);not null"`
	Address       string          `json:"address" gorm:"type:varchar(128);not null"`
	Coordinate    pq.Float64Array `json:"coordinate" gorm:"type:decimal[]"`
	Description   string          `json:"description" gorm:"type:text"`
	ImagesMain    string          `json:"imagesMain" gorm:"type:varchar(256)"`
	ImagesContent pq.StringArray  `json:"imagesContent" gorm:"type:varchar(256)[]"`
	Webpage       string          `json:"webpage" gorm:"type:varchar(32)"`
}

type LibraryCollection struct {
	gorm.Model
	ID           int                       `json:"id" gorm:"type:int;primaryKey;size:32"`
	SerialNumber string                    `json:"serialNumber" gorm:"type:varchar(32);not null"`
	LibraryID    int                       `gorm:"type:int;size:32;autoIncrement:false"`
	Library      LibraryData               `json:"library" gorm:"foreignKey:LibraryID"`
	BookID       int                       `gorm:"type:int;size:32;autoIncrement:false"`
	Book         Book                      `json:"book" gorm:"foreignKey:BookID"`
	Availability int                       `json:"availability" gorm:"not null"`
	Status       int                       `json:"status" gorm:"type:int;size:32;not null"`
	Borrow       []LibraryCollectionBorrow `json:"borrowHistory" gorm:"foreignKey:CollectionID"`
}

type LibraryCollectionBorrow struct {
	ID           int               `json:"id" gorm:"type:int;primaryKey;size:32"`
	CreatedAt    time.Time         `json:"createdAt" gorm:"type:timestamp with time zone"`
	AcceptedAt   *time.Time        `json:"acceptedAt" gorm:"type:timestamp with time zone"`
	TakedAt      *time.Time        `json:"takedAt" gorm:"type:timestamp with time zone"`
	ReturnedAt   *time.Time        `json:"returnedAt" gorm:"type:timestamp with time zone"`
	CanceledAt   *time.Time        `json:"canceledAt" gorm:"type:timestamp with time zone"`
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
	Type        string                   `json:"type" gorm:"type:varchar(16)"`
	Description datatypes.JSON           `json:"description"`
	Access      bool                     `json:"access"`
	PaperUrl    string                   `gorm:"type:varchar(128)"`
	Permission  []LibraryPaperPermission `json:"permission" gorm:"foreignKey:PaperID"`
}

type LibraryPaperPermission struct {
	gorm.Model
	ID       int            `gorm:"type:int;primaryKey;size:32"`
	PaperID  int            `gorm:"type:int;size:32;autoIncrement:false"`
	Paper    LibraryPaper   `json:"paper" gorm:"foreignKey:PaperID"`
	UserID   int            `gorm:"type:int;size:32;autoIncrement:false"`
	User     ProfileAccount `json:"user" gorm:"foreignKey:UserID"`
	Purpose  string         `json:"purpose" gorm:"type:varchar(128)"`
	Accepted *bool          `json:"access"`
}

type LibraryPaperAccess struct {
	ID           int                    `json:"id" gorm:"type:int;primaryKey;size:32"`
	CreatedAt    time.Time              `json:"createdAt" gorm:"type:timestamp with time zone"`
	PermissionID int                    `gorm:"type:int;size:32;autoIncrement:false"`
	Permission   LibraryPaperPermission `json:"paper" gorm:"foreignKey:PermissionID"`
}
