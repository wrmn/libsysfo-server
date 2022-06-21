package library

import "libsysfo-server/utility"

type LibraryData struct {
	Id          int                `json:"id"`
	Name        string             `json:"name"`
	Address     string             `json:"address"`
	Coordinate  utility.Location   `json:"coordinate"`
	Description string             `json:"description"`
	Images      utility.ImagesData `json:"images"`
}
