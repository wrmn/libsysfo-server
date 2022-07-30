package server

type borrowDataset struct {
	Month     string `json:"month"`
	Count     int64  `json:"count"`
	Requested int64  `json:"requested"`
	Taked     int64  `json:"taked"`
	Finished  int64  `json:"finished"`
	Canceled  int64  `json:"canceled"`
}

type bookDataset struct {
	Count int64 `json:"count"`
	New   int64 `json:"new"`
	Great int64 `json:"great"`
	Good  int64 `json:"good"`
	Bad   int64 `json:"bad"`
}

type paperDataset struct {
	Count   int64 `json:"count"`
	Journal int64 `json:"journal"`
	Thesis  int64 `json:"thesis"`
	Other   int64 `json:"other"`
}

type accessDataset struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type paginate struct {
	Data    int    `json:"dataTotal"`
	Next    string `json:"nextPage,omitempty"`
	Current string `json:"currentPage"`
	Prev    string `json:"previousPage"`
}

type monthCount struct {
	Borrow int64 `json:"borrow"`
	Access int64 `json:"access"`
}

type datarange struct {
	Id       int
	FromDate string
	ToDate   string
}
