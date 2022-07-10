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
	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	pg_con_string := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"))

	db, err := sql.Open("postgres", pg_con_string)
	if err != nil {
		return
	}
	DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	return
}

func Checker() {
	for {
		time.Sleep(60 * time.Minute)
		data := []LibraryCollectionBorrow{}
		DB.Find(&data)
		for _, d := range data {
			diff := time.Since(d.CreatedAt)
			if diff.Hours() >= 48 && d.Status == "requested" {
				d.Status = "finished"
				DB.Save(&d)
			}
		}
	}
}
