package database

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

func (data Book) SlugGenerator() (slug string) {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	slug = re.ReplaceAllString(strings.ToLower(data.Title), "-")
	var exist int64 = 1
	rep := 0
	for exist != 0 {
		query := DB.Preload("BookDetail").
			Where("slug = ?", slug).Find(&Book{})
		exist, _ = CheckExist(query)
		if exist != 0 {
			rep += 1
			slug = fmt.Sprintf("%s-%d", slug, rep)
		}
	}

	return
}

func CheckExist(q *gorm.DB) (int64, error) {
	return q.RowsAffected, q.Error
}
