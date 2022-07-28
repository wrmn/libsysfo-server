package database

import "gorm.io/gorm"

type Feedback struct {
	gorm.Model
	Name    string  `json:"name" gorm:"type:varchar(32);not null"`
	Email   *string `json:"email" gorm:"type:varchar(64)"`
	Message string  `json:"message"`
}

type ThirdPartyJobs struct {
	gorm.Model
	IssuerID     string `json:"issuer" gorm:"varchar(128)"`
	Job          string `json:"job" gorm:"type:varchar(32)"`
	Destination  string `json:"destination" gorm:"varchar(32)"`
	ResponseBody string `json:"responseBody" `
	Status       int    `json:"status" gorm:"type:int;size:32;autoIncrement:false"`
}

type bookDataset struct {
	Books []Book `json:"books"`
}

type paperDataset struct {
	Papers []LibraryPaper `json:"papers"`
}
