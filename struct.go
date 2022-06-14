package main

type ImagesData struct {
	Main    string   `json:"main"`
	Content []string `json:"content"`
}

type NameData struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	Coordinate  Location   `json:"coordinate"`
	Description string     `json:"description"`
	Images      ImagesData `json:"images"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type NameDatas struct {
	Data []NameData `json:"data"`
}
