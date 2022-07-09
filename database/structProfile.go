package database

import (
	"time"

	"gorm.io/gorm"
)

// login data
type ProfileAccount struct {
	gorm.Model
	ID          int         `json:"id" gorm:"type:int;primaryKey;size:32"`
	Username    *string     `json:"username" gorm:"type:varchar(32);unique"`
	Email       string      `json:"email" gorm:"type:varchar(64);unique;not null"`
	AccountType int         `json:"accountType" gorm:"type:int;size:32"`
	Password    string      `json:"password,omitempty" gorm:"type:char(32);not null"`
	LastLogin   *time.Time  `json:"lastLogin" gorm:"type:timestamp"`
	ProfileData ProfileData `json:"profileData" gorm:"foreignKey:UserID"`
}

// account information
type ProfileData struct {
	UserID       int        `json:"userId" gorm:"type:int;size:32;primaryKey;autoIncrement:false"`
	VerifiedAt   *time.Time `json:"verifiedAt" gorm:"type:timestamp"`
	Name         string     `json:"name" gorm:"type:varchar(32);not null"`
	Gender       *string    `json:"gender" gorm:"type:char(1)"`
	PlaceOfBirth *string    `json:"placeOfBirth" gorm:"type:varchar(32)"`
	DateOfBirth  *string    `json:"dateOfBirth" gorm:"type:date"`
	Address1     *string    `json:"address" gorm:"type:varchar(128)"`
	Profession   *string    `json:"profession" gorm:"type:varchar(32)"`
	Institution  *string    `json:"institution" gorm:"type:varchar(64)"`
	PhoneCode    *string    `json:"phoneCode" gorm:"type:varchar(5)"`
	PhoneNo      *string    `json:"phoneNo" gorm:"type:varchar(20)"`
	IsWhatsapp   bool       `json:"isWhatsapp" gorm:"not null"`
	Images       string     `json:"images" gorm:"not null;type:varchar(128)"`
}
