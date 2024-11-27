package ui

import "github.com/maddalax/htmgo/framework/h"

type Item struct {
	Value string
	Text  string
}

type SelectProps struct {
	Id       string
	Name     string
	Required bool
	Value    string
	Children []h.Ren
	Items    []Item
}

func Select(props SelectProps) *h.Element {
	return h.Select(
		h.If(
			props.Required,
			h.Required(),
		),
		h.Class("w-full rounded border p-2 focus:outline-none focus:ring-0 focus:border-gray-400"),
		h.If(
			props.Id != "",
			h.Id(props.Id),
		),
		h.If(
			props.Name != "",
			h.Name(props.Name),
		),
		h.List(props.Items, func(item Item, index int) *h.Element {
			return h.Option(
				h.If(
					item.Value == props.Value,
					h.Selected(),
				),
				h.Value(item.Value),
				h.Text(item.Text),
			)
		}),
	)
}
