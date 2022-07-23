package database

import (
	"gorm.io/gorm"
)

type Book struct {
	ID         int        `json:"id,omitempty" gorm:"type:int;primaryKey;size:32"`
	BookDetail BookDetail `json:"detail" gorm:"foreignKey:ID"`
	Image      string     `json:"image" gorm:"type:varchar(256)"`
	Title      string     `json:"title" gorm:"type:varchar(256)"`
	Author     string     `json:"author" gorm:"type:varchar(256)"`
	Source     string     `json:"source" gorm:"type:varchar(32)"`
	Slug       string     `json:"slug" gorm:"type:varchar(256)"`
}

type BookDetail struct {
	gorm.Model
	ID          int    `json:"id,omitempty" gorm:"type:int;primaryKey;size:32;autoIncrement:false"`
	ReleaseDate string `json:"release_date" gorm:"type:varchar(32)"`
	Description string `json:"description" gorm:"type:text"`
	Language    string `json:"language" gorm:"type:varchar(64)"`
	Country     string `json:"country" gorm:"type:varchar(32)"`
	Publisher   string `json:"publisher" gorm:"type:varchar(32)"`
	PageCount   int    `json:"page_count" gorm:"type:int;size:32"`
	Category    string `json:"category" gorm:"type:varchar(64)"`
}
