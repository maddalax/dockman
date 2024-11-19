package ui

import (
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
	"time"
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
	Post           string
	Trigger        string
	Delay          time.Duration
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
		h.OnLoad(
			// language=JavaScript
			h.EvalJs(`
		  let form = self.closest('form');
      let startLoader = new Function(self.dataset.startLoader);
			if(form) {
       form.addEventListener('submit', function() {
         startLoader.call(self);
       })
      } else {
         // if the button is not in a form, we need to manually trigger the click event
	       self.addEventListener('click', function() {
           startLoader.call(self);
         })
      }
     `),
		),
		h.OnEvent("data-start-loader",
			js.RemoveClassOnChildren(".loading", "hidden"),
			js.SetClassOnChildren(".submit", "hidden"),
		),
		h.HxAfterRequest(
			// delay so the loading spinner doesn't flash too quickly
			// and we give some feedback to the user
			js.RunAfterTimeout(props.Delay,
				js.SetClassOnChildren(".loading", "hidden"),
				js.RemoveClassOnChildren(".submit", "hidden"),
			),
		),
		h.Class("flex gap-2 justify-center"),

		h.Div(
			h.OnLoad(
				// let's make sure the button is the same width as the loading spinner
				js.EvalJs(`
					const button = self.nextElementSibling.getBoundingClientRect();
					self.style.width = button.width + 'px';
				`),
			),
			h.Class("loading hidden text-center", buttonClasses),
			h.Disabled(),
			h.Div(
				h.Class("flex gap-2 items-center justify-start"),
				spinner(),
				h.Text(h.Ternary(props.SubmittingText != "", props.SubmittingText, "Loading...")),
			),
		),

		h.Button(
			h.Type("submit"),
			h.Class("submit", buttonClasses),
			h.Text(props.Text),
			h.If(
				props.Post != "",
				h.Post(props.Post),
			),
		),
	)
}

func spinner(children ...h.Ren) *h.Element {
	return h.Div(
		h.Children(children...),
		h.Class("spinner spinner-border animate-spin inline-block w-4 h-4 border-2 rounded-full border-slate-200 border-t-transparent"),
		h.Attribute("role", "status"),
	)
}
