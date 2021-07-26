package doc

import (
	"fmt"
)

// Table ...
type Table struct {
	columns   []column
	columnMap map[string]column
}

// Rows ...
type Rows struct {
	table *Table
}

type column struct {
	Title string
	Width int
}

// Column ...
func (table *Table) Column(title string, width int) *Table {
	existing, found := table.columnMap[title]
	if found {
		existing.Width = width
		return table
	}
	column := column{Title: title, Width: width}
	table.columns = append(table.columns, column)
	table.columnMap[title] = column
	return table
}

// StartRows ...
func (table *Table) StartRows() *Rows {

	totalLength := 0
	for _, col := range table.columns {
		Cell(FitLeft(col.Title, col.Width))
		totalLength += col.Width
	}
	fmt.Fprintln(writer())
	Divider(totalLength)
	return &Rows{table: table}
}

// Row Print Row
func (rows *Rows) Row(params ...string) *Rows {
	for index, rowValue := range params {
		col := rows.table.columns[index]
		Cell(FitLeft(rowValue, col.Width))
	}
	fmt.Fprintln(writer())
	return rows
}

// Divider Print Row
func (rows *Rows) Divider() *Rows {
	totalLength := 0
	for _, col := range rows.table.columns {
		totalLength += col.Width
	}
	Divider(totalLength)
	return rows
}

// NewTable ...
func NewTable() *Table {
	mp := make(map[string]column, 0)
	ls := make([]column, 0)
	return &Table{columnMap: mp, columns: ls}
}
