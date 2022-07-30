package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() (err error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// pg_con_string := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_USERNAME"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_DATABASE"))

	// db, err := sql.Open("postgres", pg_con_string)
	if err != nil {
		return
	}
	DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	return
}

func Checker() {
	var err error
	for {
		time.Sleep(60 * time.Minute)
		now := time.Now()

		dataBorrow := []LibraryCollectionBorrow{}
		dataAccess := []LibraryPaperPermission{}

		err = DB.Find(&dataAccess).Error
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for _, d := range dataBorrow {
			if d.TakedAt == nil {
				var diff time.Duration
				if d.AcceptedAt == nil {
					diff = time.Since(d.CreatedAt)
				} else {
					diff = time.Since(*d.AcceptedAt)
				}
				if diff.Hours() >= 48 {
					d.CanceledAt = &now
					DB.Save(&d)
				}
			}
		}

		err = DB.Find(&dataBorrow).Error
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		for _, d := range dataAccess {
			if d.AcceptedAt == nil && d.CanceledAt == nil {
				diff := time.Since(d.CreatedAt)
				if diff.Hours() <= 48 {
					d.CanceledAt = &now
					DB.Save(&d)
				}
			}
		}
	}
}
