package book

import "gorm.io/gorm"

type BookData struct {
	gorm.Model
	Id        int    `json:"id" gorm:"type:int;primaryKey;size:32"`
	Serial    int    `json:"serial" gorm:"type:int;size:32;not null"`
	Title     string `json:"title" gorm:"type:varchar(96);not null"`
	Author    string `json:"author" gorm:"type:varchar(32);not null"`
	Publisher string `json:"publisher" gorm:"type:varchar(32);not null"`
}
