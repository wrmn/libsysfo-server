package database

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Book struct {
	Id         int        `json:"id,omitempty" gorm:"type:int;primaryKey;size:32"`
	BookDetail BookDetail `json:"detail" gorm:"foreignKey:BookId"`
	Image      string     `json:"image" gorm:"type:varchar(256)"`
	Title      string     `json:"title" gorm:"type:varchar(128)"`
	Author     string     `json:"author" gorm:"type:varchar(64)"`
	Slug       string     `json:"slug" gorm:"type:varchar(128)"`
}

type BookDetail struct {
	gorm.Model
	Id          int            `json:"id,omitempty" gorm:"type:int;primaryKey;size:32"`
	BookId      int            `json:"bookId,omitempty" gorm:"type:int;size:32"`
	ReleaseDate datatypes.Date `json:"release_date" gorm:"type:date"`
	Description string         `json:"description" gorm:"type:text"`
	Language    string         `json:"language" gorm:"type:varchar(64)"`
	Country     string         `json:"country" gorm:"type:varchar(32)"`
	Publisher   string         `json:"publisher" gorm:"type:varchar(32)"`
	PageCount   int            `json:"page_count" gorm:"type:int;size:32"`
	Category    string         `json:"category" gorm:"type:varchar(64)"`
}
