package database

import (
	"libsysfo-server/database/book"
	"libsysfo-server/database/library"
)

type LibrariesData struct {
	Data []library.LibraryData `json:"data"`
}

type BooksData struct {
	Data []book.BookData `json:"data"`
}
