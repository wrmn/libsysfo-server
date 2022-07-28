package database

import (
	"database/sql"
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
	for {
		time.Sleep(60 * time.Minute)
		dataBorrow := []LibraryCollectionBorrow{}
		DB.Find(&dataBorrow)
		for _, d := range dataBorrow {
			var diff time.Duration
			if d.AcceptedAt == nil {
				diff = time.Since(d.CreatedAt)
			} else {
				diff = time.Since(*d.AcceptedAt)
			}
			if diff.Hours() >= 48 && d.TakedAt == nil {
				now := time.Now()
				d.CanceledAt = &now
				DB.Save(&d)
			}
		}
		dataAccess := []LibraryPaperPermission{}
		DB.Find(&dataAccess)
		stats := false
		for _, d := range dataAccess {
			diff := time.Since(d.CreatedAt)
			if diff.Hours() <= 48 && d.Accepted == nil {
				d.Accepted = &stats
				DB.Save(&d)
			}
		}
	}
}
