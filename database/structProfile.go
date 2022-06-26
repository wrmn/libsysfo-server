package database

import (
	"time"

	"gorm.io/gorm"
)

// login data
type ProfileAccount struct {
	gorm.Model
	Id        int        `json:"id" gorm:"type:int;primaryKey;size:32"`
	Username  string     `json:"username" gorm:"type:varchar(24);unique;not null"`
	Email     string     `json:"email" gorm:"type:varchar(64);unique;not null"`
	Password  string     `json:"password" gorm:"type:char(32);not null"`
	LastLogin *time.Time `json:"lastLogin" gorm:"type:timestamp"`
}

// account information
type ProfileData struct {
	UserId       int            `json:"userId" gorm:"type:int;size:32;primaryKey;autoIncrement:false"`
	User         ProfileAccount `json:"user" gorm:"foreignKey:UserId"`
	VerifiedAt   *time.Time     `json:"verifiedAt" gorm:"type:timestamp"`
	Name         string         `json:"name" gorm:"type:varchar(32);not null"`
	Gender       string         `json:"gender" gorm:"type:char(1);not null"`
	PlaceOfBirth string         `json:"placeOfBirth" gorm:"type:varchar(32);not null"`
	DateOfBirth  string         `json:"dateOfBirth" gorm:"type:date;not null"`
	Address1     string         `json:"address" gorm:"type:varchar(128);not null"`
	Profession   string         `json:"profession" gorm:"type:varchar(32);not null"`
	Institution  string         `json:"institution" gorm:"type:varchar(64)not null"`
	PhoneNo      string         `json:"phoneNo" gorm:"type:varchar(20);not null"`
	IsWhatsapp   bool           `json:"isWhatsapp" gorm:"not null"`
	Images       string         `json:"images" gorm:"not null;type:varchar(128)"`
}
