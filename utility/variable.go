package utility

var (
	Dmy    string = "2006-01-02"
	Dmyhms string = "15:04:05.000"
)

type ImagesData struct {
	Main    string   `json:"main"`
	Content []string `json:"content"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
