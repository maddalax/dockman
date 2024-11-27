package ui

import (
	"github.com/maddalax/htmgo/framework/h"
)

type ChoiceCardProps struct {
	Id             string
	Icon           *h.Element
	Title          string
	InputName      string
	InputValue     string
	DefaultChecked bool
	Description    string
	InputProps     *h.AttributeMapOrdered
}

func ChoiceCard(props ChoiceCardProps) *h.Element {
	return h.Label(
		h.For(props.Id),
		h.ClassX("w-full flex flex-col rounded-lg text-center mr-4 mb-4 shadow-sm border-2 border-slate-100", h.ClassMap{
			"bg-white has-[:checked]:bg-brand-50 has-[:checked]:border-brand-200": true,
		}),
		h.Div(
			h.Class("flex flex-1 flex-col p-4"),
			h.If(
				props.Icon != nil,
				h.Div(
					h.Class("mx-auto max-h-16 max-w-16 shrink-0 rounded-full"),
					props.Icon,
				),
			),
			h.H3(
				h.ClassX("text-sm font-medium text-gray-900", h.ClassMap{
					"mt-4": props.Icon != nil,
				}),
				h.Text(props.Title),
			),
			h.Dl(
				h.Class("mt-1 flex flex-col justify-between"),
				h.Dt(
					h.Class("sr-only"),
					h.Text("Title"),
				),
				h.Dd(
					h.Class("text-sm text-gray-500"),
					h.Text(props.Description),
				),
				h.Dt(
					h.Class("sr-only"),
					h.Text("Role"),
				),
			),
		),
		h.Input(
			"radio",
			h.Name(props.InputName),
			h.Value(props.InputValue),
			h.Id(props.Id),
			h.If(
				props.DefaultChecked,
				h.Checked(),
			),
			h.If(props.InputProps != nil, props.InputProps),
			h.Class("sr-only"),
		),
	)
}
