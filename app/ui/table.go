package ui

import (
	"github.com/maddalax/htmgo/framework/h"
)

type Table struct {
	columns    []*h.Element
	rows       [][]*h.Element
	currentRow int
}

func NewTable() *Table {
	return &Table{
		columns: make([]*h.Element, 0),
		rows:    make([][]*h.Element, 0),
	}
}

func (t *Table) AddColumn(column string) {
	t.columns = append(t.columns, h.Th(
		h.Class("py-2 px-4 text-sm font-semibold text-gray-700"),
		h.Text(column),
	))
}

func (t *Table) AddColumns(columns []string) {
	for _, column := range columns {
		t.AddColumn(column)
	}
}

func (t *Table) AddRow() {
	t.rows = append(t.rows, make([]*h.Element, 0))
	t.currentRow = len(t.rows) - 1
}

func (t *Table) AddCell(cell *h.Element) {
	t.rows[t.currentRow] = append(t.rows[t.currentRow], cell)
}

func (t *Table) AddCellText(cell string) {
	t.rows[t.currentRow] = append(t.rows[t.currentRow], h.Pf(cell))
}

func (t *Table) WithCells(cells ...*h.Element) {
	for _, cell := range cells {
		t.AddCell(cell)
	}
}

func (t *Table) WithCellTexts(cells ...string) {
	for _, cell := range cells {
		t.AddCellText(cell)
	}
}

func (t *Table) Render() *h.Element {
	return h.Table(
		h.Class("w-full border-collapse border border-gray-300 x-overflow-auto truncate"),
		h.THead(
			h.Tr(
				h.Class("bg-gray-100 text-left border-b border-gray-300"),
				h.List(t.columns, func(column *h.Element, index int) *h.Element {
					return column
				}),
			),
		),
		h.TBody(
			h.List(t.rows, func(row []*h.Element, index int) *h.Element {
				return h.Tr(
					h.Class("border-b border-gray-300 hover:bg-gray-50"),
					h.List(row, func(cell *h.Element, index int) *h.Element {
						return h.Td(
							h.Class("py-2 px-4 text-sm text-gray-700"),
							cell,
						)
					},
					),
				)
			}),
		),
	)
}
