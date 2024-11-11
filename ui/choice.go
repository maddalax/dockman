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
		h.ClassX("flex flex-col rounded-lg text-center w-[250px] max-w-[250px] mr-4 mb-4 shadow-sm border-2 border-slate-100", h.ClassMap{
			"bg-white has-[:checked]:bg-brand-50 has-[:checked]:border-brand-200": true,
		}),
		h.Div(
			h.Class("flex flex-1 flex-col p-4"),
			h.Div(
				h.Class("mx-auto h-32 w-32 shrink-0 rounded-full"),
				props.Icon,
			),
			h.H3(
				h.Class("mt-6 text-sm font-medium text-gray-900"),
				h.Text(props.Title),
			),
			h.Dl(
				h.Class("mt-1 flex grow flex-col justify-between"),
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
			props.InputProps,
			h.Class("sr-only"),
		),
	)
}
