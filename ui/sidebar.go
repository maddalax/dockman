package ui

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/resources"
	"paas/urls"
)

type Section struct {
	Title string
	Pages []*Page
}

type Page struct {
	Title string
	Path  string
}

func DocPath(path string) string {
	return "/docs" + path
}

func MainSidebar(ctx *h.RequestContext) *h.Element {
	return h.Div(
		h.Class("px-3 py-2 pr-6 md:min-h-screen pb-4 mb:pb-0 bg-neutral-50 border-r border-r-slate-300 overflow-y-auto"),
		h.Div(
			h.Class("flex flex-col gap-4"),
			h.Div(
				h.Class("mb-3"),
				h.A(
					h.Href("/"),
					h.Text("paas"),
					h.Class("md:mt-4 text-xl text-slate-900 font-bold"),
				),
			),
			ResourceList(ctx),
			DebugSection(),
		),
	)
}

func ResourceList(ctx *h.RequestContext) *h.Element {
	names := resources.GetNames(ctx.ServiceLocator())

	return h.Div(
		h.Class("flex flex-col gap-2"),
		h.Div(
			h.Class("flex justify-between items-center"),
			h.P(
				h.Text("Resources"),
				h.Class("text-slate-800 font-bold"),
			),
			h.A(
				h.Text("+ New"),
				h.Href(urls.NewResourceUrl()),
				h.Class("bg-slate-900 hover:bg-slate-800 text-white text-xs font-bold py-2 px-2 rounded"),
			),
		),
		h.Div(
			h.Class("flex flex-col gap-2"),
			h.List(names, func(resource resources.ResourceName, index int) *h.Element {
				return h.A(
					h.Href(urls.ResourceUrl(resource.Id)),
					h.Text(resource.Name),
					h.Class("text-slate-900 hover:text-brand-400"),
				)
			}),
		),
	)
}

func DebugSection() *h.Element {

	links := []Page{
		{
			Title: "KV Viewer",
			Path:  "/debug/jetstream/kv",
		},
		{
			Title: "Stream Viewer",
			Path:  "/debug/jetstream/streams",
		},
	}

	return h.Div(
		h.Class("flex flex-col gap-2"),
		h.Div(
			h.Class("flex justify-between items-center"),
			h.P(
				h.Text("Debug"),
				h.Class("text-slate-800 font-bold"),
			),
		),
		h.Div(
			h.Class("flex flex-col gap-2"),
			h.List(links, func(link Page, index int) *h.Element {
				return h.A(
					h.Href(link.Path),
					h.Text(link.Title),
					h.Class("text-slate-900 hover:text-brand-400"),
				)
			}),
		),
	)
}
