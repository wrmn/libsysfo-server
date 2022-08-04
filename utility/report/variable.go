package report

type Table struct {
	Header []MainHeader
	Table  []MainTable
	Data   int
}

type MainTable struct {
	Width float64
	Name  string
}
type MainHeader struct {
	Name  string
	Value interface{}
}

var BookReport []MainTable = []MainTable{{
	Width: 5,
	Name:  "No",
}, {
	Width: 20,
	Name:  "Created At",
}, {
	Width: 20,
	Name:  "Serial Number",
}, {
	Width: 40,
	Name:  "Title",
}, {
	Width: 30,
	Name:  "Category",
}, {
	Width: 20,
	Name:  "Author",
}, {
	Width: 20,
	Name:  "Release Date",
}, {
	Width: 20,
	Name:  "Publisher",
}, {
	Width: 20,
	Name:  "Language",
}, {
	Width: 20,
	Name:  "Country",
}, {
	Width: 10,
	Name:  "Page",
}, {
	Width: 20,
	Name:  "Availability",
}, {
	Width: 20,
	Name:  "Status",
}, {
	Width: 10,
	Name:  "Borrow Total",
}}

var BorrowReport []MainTable = []MainTable{
	{
		Width: 5,
		Name:  "No",
	},
	{
		Width: 15,
		Name:  "Status",
	},
	{
		Width: 20,
		Name:  "Requested At",
	},
	{
		Width: 20,
		Name:  "Accepted At",
	},
	{
		Width: 20,
		Name:  "Taked At",
	},
	{
		Width: 20,
		Name:  "Returned At",
	},
	{
		Width: 20,
		Name:  "Canceled At",
	},
	{
		Width: 20,
		Name:  "Book Title",
	},
	{
		Width: 20,
		Name:  "Book Serial Number",
	},
	{
		Width: 20,
		Name:  "Borrower Name",
	},
	{
		Width: 20,
		Name:  "Borrower Username",
	},
	{
		Width: 20,
		Name:  "Borrower E-mail",
	}}
