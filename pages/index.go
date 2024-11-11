package pages

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/ui"
)

func Index(ctx *h.RequestContext) *h.Page {
	return SidebarPage(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-2"),
			Title("Introduction"),
			Text(`
				htmgo is a lightweight pure go way to build interactive websites / web applications using go & htmx.
				We give you the utilities to build html using pure go code in a reusable way (go functions are components) while also providing htmx functions to add interactivity to your app.
			`),
			h.P(
				Link("The site you are reading now", "https://github.com/maddalax/htmgo/tree/master/htmgo-site"),
				h.Text(" was written with htmgo!"),
			),
			NextStep(
				"mt-4",
				h.Div(),
				NextBlock("Getting Started", ui.DocPath("/installation")),
			),
		),
	)
}
