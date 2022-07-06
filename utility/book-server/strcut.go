package bookserver

type BookResponse struct {
	Book   *Book   `json:"book,omitempty"`
	Books  []Book  `json:"books,omitempty"`
	Links  *Links  `json:"links,omitempty"`
	Meta   *Meta   `json:"meta,omitempty"`
	Status *int64  `json:"status,omitempty"`
	Type   *string `json:"type,omitempty"`
	Title  *string `json:"title,omitempty"`
}

type Book struct {
	Image       *string `json:"image,omitempty"`
	Title       *string `json:"title,omitempty"`
	Author      *string `json:"author,omitempty"`
	Price       *string `json:"price,omitempty"`
	OriginalURL *string `json:"original_url,omitempty"`
	URL         *string `json:"url,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	Detail      *Detail `json:"detail,omitempty"`
}

type Detail struct {
	ReleaseDate *string  `json:"release_date,omitempty"`
	Description *string  `json:"description,omitempty"`
	Language    *string  `json:"language,omitempty"`
	Country     *string  `json:"country,omitempty"`
	Publisher   *string  `json:"publisher,omitempty"`
	PageCount   *float64 `json:"page_count,omitempty"`
	Category    *string  `json:"category,omitempty"`
}

type Links struct {
	First *string     `json:"first,omitempty"`
	Last  interface{} `json:"last"`
	Prev  interface{} `json:"prev"`
	Next  *string     `json:"next,omitempty"`
}

type Meta struct {
	CurrentPage *int64  `json:"current_page,omitempty"`
	From        *int64  `json:"from,omitempty"`
	Path        *string `json:"path,omitempty"`
	PerPage     *int64  `json:"per_page,omitempty"`
	To          *int64  `json:"to,omitempty"`
}
