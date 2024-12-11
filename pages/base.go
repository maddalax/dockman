package pages

import (
	"dockman/app"
	"dockman/app/ui"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
	"strings"
)

func Title(title string) *h.Element {
	return h.H1(
		h.Text(title),
		h.Class("text-2xl font-bold"),
	)
}

func NextStep(classes string, prev *h.Element, next *h.Element) *h.Element {
	return h.Div(
		h.Class("flex gap-2 justify-between", classes),
		prev,
		next,
	)
}

func NextBlock(text string, url string) *h.Element {
	return h.A(
		h.Href(url),
		h.Class("w-[50%] border border-slate-300 p-4 rounded text-right hover:border-blue-400 cursor-pointer"),
		h.P(
			h.Text("Next"),
			h.Class("text-slate-600 text-sm"),
		),
		h.P(
			h.Text(text),
			h.Class("text-blue-500 hover:text-blue-400"),
		),
	)
}

func Text(text string) *h.Element {
	split := strings.Split(text, "\n")
	return h.Div(
		h.Class("flex flex-col gap-2 leading-relaxed text-slate-900 break-words"),
		h.List(split, func(item string, index int) *h.Element {
			return h.P(
				h.UnsafeRaw(item),
			)
		}),
	)
}

func Link(text string, href string, additionalClasses ...string) *h.Element {
	additionalClasses = append(additionalClasses, "text-blue-500 hover:text-blue-400")
	return h.A(
		h.Href(href),
		h.Text(text),
		h.Class(
			additionalClasses...,
		),
	)
}

func NavBar(ctx *h.RequestContext) *h.Element {
	var OpenMobileSidebarButton = h.Button(
		h.OnClick(
			js.EvalCommandsOnSelector(
				"#mobile-sidebar",
				js.AddClass("relative"),
				js.RemoveClass("hidden"),
			),
		),
		h.Type("button"),
		h.Class("-m-2.5 p-2.5 text-gray-700 lg:hidden"),
		h.Span(
			h.Class("sr-only"),
			h.Text("Open sidebar"),
		),
		h.Svg(
			h.Class("size-6"),
			h.Attribute("fill", "none"),
			h.Attribute("viewBox", "0 0 24 24"),
			h.Attribute("stroke-width", "1.5"),
			h.Attribute("stroke", "currentColor"),
			h.Attribute("aria-hidden", "true"),
			h.Attribute("data-slot", "icon"),
			h.Path(
				h.Attribute("stroke-linecap", "round"),
				h.Attribute("stroke-linejoin", "round"),
				h.Attribute("d", "M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"),
			),
		),
	)

	return h.Div(
		h.Class("sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b border-gray-200 bg-white px-4 shadow-sm sm:gap-x-6 sm:px-6 lg:px-8"),
		OpenMobileSidebarButton,
		h.Div(
			h.Class("h-6 w-px bg-gray-200 lg:hidden"),
			h.Attribute("aria-hidden", "true"),
		),
		SearchBar(ctx),
	)
}

func SearchBar(ctx *h.RequestContext) *h.Element {
	user := app.CurrentUser(ctx)

	return h.Div(
		h.Class("flex flex-1 gap-x-4 self-stretch lg:gap-x-6"),
		h.Form(
			h.Class("relative flex flex-1"),
			h.Action("#"),
			h.Method("GET"),
			h.Label(
				h.For("search-field"),
				h.Class("sr-only"),
				h.Text("Search"),
			),
			h.Svg(
				h.Class("pointer-events-none absolute inset-y-0 left-0 h-full w-5 text-gray-400"),
				h.Attribute("viewBox", "0 0 20 20"),
				h.Attribute("fill", "currentColor"),
				h.Attribute("aria-hidden", "true"),
				h.Attribute("data-slot", "icon"),
				h.Path(
					h.Attribute("fill-rule", "evenodd"),
					h.Attribute("d", "M9 3.5a5.5 5.5 0 1 0 0 11 5.5 5.5 0 0 0 0-11ZM2 9a7 7 0 1 1 12.452 4.391l3.328 3.329a.75.75 0 1 1-1.06 1.06l-3.329-3.328A7 7 0 0 1 2 9Z"),
					h.Attribute("clip-rule", "evenodd"),
				),
			),
			h.Input(
				"search",
				h.Id("search-field"),
				h.Class("block size-full border-0 py-0 pl-8 pr-0 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm"),
				h.Placeholder("Search..."),
				h.Name("search"),
			),
		),
		h.Div(
			h.Class("flex gap-2 items-center"),
			h.If(user != nil, h.Div(
				h.Class("flex-shrink-0"),
				h.Pf("%s", user.Email),
			)),
		),
	)
}

func SidebarPage(ctx *h.RequestContext, children ...h.Ren) *h.Page {
	return RootPage(
		ctx,
		h.Div(
			h.Div(
				ui.MainSidebar(ctx),
				h.Div(
					h.Class("lg:pl-60"),
					NavBar(ctx),
					h.Main(
						h.Class("py-10"),
						h.Div(
							h.Class("px-4 sm:px-6 lg:px-8"),
							h.Children(children...),
						),
					),
				),
			),
		),
	)
}
