package report

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func Cell(s string, i int) string {
	return fmt.Sprintf("%s%d", s, i)
}

func getColumnName(col int) string {
	name := make([]byte, 0, 3)
	const aLen = 'Z' - 'A' + 1
	for ; col > 0; col /= aLen + 1 {
		name = append(name, byte('A'+(col-1)%aLen))
	}
	for i, j := 0, len(name)-1; i < j; i, j = i+1, j-1 {
		name[i], name[j] = name[j], name[i]
	}
	return string(name)
}

// b []MainHeader, m []MainTable
func (t Table) CreateMainTable(sheet string) *excelize.File {
	f := excelize.NewFile()
	for i, k := range t.Header {
		f.SetCellValue(sheet, Cell("A", i+1), k.Name)
		f.SetCellValue(sheet, Cell("C", i+1), k.Value)
	}

	for i, c := range t.Table {
		f.SetCellValue(sheet, Cell(getColumnName(i+1), len(t.Header)+1), c.Name)
		f.SetColWidth(sheet, getColumnName(i+1), getColumnName(i+1), c.Width)
	}
	return f
}

// b for general data on head : Head
// d for total data : Data
// m for total header : Header
func (c Table) Styling(sheet string, f *excelize.File) {
	header, _ := f.NewStyle(`{"border": [{ "type": "left", "color": "000000", "style": 5 },{ "type": "top", "color": "000000", "style": 5 },{ "type": "bottom", "color": "000000", "style": 5 },{ "type": "right", "color": "000000", "style": 5 }]}`)
	content, _ := f.NewStyle(`{"border": [{ "type": "left", "color": "000000", "style": 5 },{ "type": "right", "color": "000000", "style": 5 }], "alignment":{"wrap_text":true, "horizontal":"center", "vertical":"center"}}`)
	footer, _ := f.NewStyle(`{"border": [{ "type": "left", "color": "000000", "style": 5 },{ "type": "right", "color": "000000", "style": 5 },{ "type": "bottom", "color": "000000", "style": 5 }], "alignment":{"wrap_text":true, "horizontal":"center", "vertical":"center"}}`)

	f.SetCellStyle(sheet, Cell("A", len(c.Header)+1), Cell(getColumnName(len(c.Table)), len(c.Header)+1), header)
	f.SetCellStyle(sheet, Cell("A", len(c.Header)+2), Cell(getColumnName(len(c.Table)), len(c.Header)+c.Data), content)
	f.SetCellStyle(sheet, Cell("A", len(c.Header)+c.Data+1), Cell(getColumnName(len(c.Table)), len(c.Header)+c.Data+1), footer)
}
