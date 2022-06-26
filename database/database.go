package database

import (
	"database/sql"
	"os"

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
