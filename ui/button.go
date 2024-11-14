package ui

import (
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
)

type ButtonProps struct {
	Id       string
	Text     string
	Target   string
	Type     string
	Trigger  string
	Get      string
	Post     string
	Href     string
	Class    string
	Children []h.Ren
}

type SubmitButtonProps struct {
	Text           string
	SubmittingText string
	Class          string
}

func PrimaryButton(props ButtonProps) *h.Element {
	props.Class = h.MergeClasses(props.Class, "border-slate-800 bg-slate-900 hover:bg-slate-800 text-white")
	return Button(props)
}

func SecondaryButton(props ButtonProps) *h.Element {
	props.Class = h.MergeClasses(props.Class, "border-gray-700 bg-gray-700 text-white")
	return Button(props)
}

func Button(props ButtonProps) *h.Element {

	text := h.Text(props.Text)

	tag := h.Ternary(props.Href != "", "a", "button")

	button := h.Tag(tag,
		h.If(
			props.Id != "",
			h.Id(props.Id),
		),
		h.If(
			props.Href != "",
			h.Href(props.Href),
		),
		h.If(
			props.Children != nil,
			h.Children(props.Children...),
		),
		h.Class("flex gap-1 items-center justify-center border p-2 rounded cursor-hover", props.Class),
		h.If(
			props.Get != "",
			h.Get(props.Get),
		),
		h.If(
			props.Post != "",
			h.Post(props.Post),
		),
		h.If(
			props.Target != "",
			h.HxTarget(props.Target),
		),
		h.IfElse(
			props.Type != "",
			h.Type(props.Type),
			h.Type("button"),
		),
		text,
	)

	return button
}

func SubmitButton(props SubmitButtonProps) *h.Element {
	buttonClasses := h.MergeClasses(
		"rounded items-center px-3 py-2 border-slate-800 bg-slate-900 hover:bg-slate-800 text-white w-full text-center", props.Class)

	return h.Div(
		h.HxBeforeRequest(
			js.RemoveClassOnChildren(".loading", "hidden"),
			js.SetClassOnChildren(".submit", "hidden"),
		),
		h.HxAfterRequest(
			js.SetClassOnChildren(".loading", "hidden"),
			js.RemoveClassOnChildren(".submit", "hidden"),
		),
		h.Class("flex gap-2 justify-center"),
		h.Button(
			h.Class("loading hidden relative text-center", buttonClasses),
			spinner(),
			h.Disabled(),
			h.Text(props.SubmittingText),
		),
		h.Button(
			h.Type("submit"),
			h.Class("submit", buttonClasses),
			h.Text(props.Text),
		),
	)
}

func spinner(children ...h.Ren) *h.Element {
	return h.Div(
		h.Children(children...),
		h.Class("absolute left-1 spinner spinner-border animate-spin inline-block w-6 h-6 border-4 rounded-full border-slate-200 border-t-transparent"),
		h.Attribute("role", "status"),
	)
}
