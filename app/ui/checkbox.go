package ui

import "github.com/maddalax/htmgo/framework/h"

type CheckboxProps struct {
	Id      string
	Label   string
	Name    string
	Checked bool
}

func Checkbox(props CheckboxProps) *h.Element {
	return h.Label(
		h.For(props.Id),
		h.Class("flex cursor-pointer items-start gap-4"),
		h.Div(
			h.Class("flex items-center"),
			h.Text("​"),
			h.Input(
				"checkbox",
				h.Class("size-4 rounded border-gray-300"),
				h.Id(props.Id),
				h.Ternary(props.Checked, h.Checked(), nil),
			),
		),
		h.Div(
			h.Strong(
				h.Class("font-medium text-gray-900"),
				h.Text(props.Label),
			),
		),
	)
}
