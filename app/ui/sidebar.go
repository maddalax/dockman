package ui

import (
	"dockside/app"
	"dockside/app/urls"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
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
	return h.Fragment(
		h.Div(
			h.Class("lg:hidden"),
			MobileSidebar(ctx),
		),
		DesktopSidebar(ctx),
	)
}

func DesktopSidebar(ctx *h.RequestContext) *h.Element {
	return h.Div(
		h.Class("hidden lg:fixed lg:inset-y-0 lg:z-40 lg:flex lg:w-60 lg:flex-col"),
		h.Div(
			h.Class("fixed inset-y-0 z-40 flex w-60 flex-col"),
			h.Div(
				h.Class("flex grow flex-col gap-y-5 overflow-y-auto border-r border-gray-200 bg-white px-4 pb-4"),
				h.Div(
					h.Class("mt-2 flex h-10 shrink-0 items-center"),
					h.Div(
						h.Class("h-6 w-auto"),
						h.H3F("dockside", h.Class("text-xl font-bold")),
					),
				),
				RoutingSection(),
				ResourceList(ctx),
				ServerList(ctx),
				DebugSection(),
			),
		),
	)
}

func MobileSidebar(ctx *h.RequestContext) *h.Element {
	CloseButton := h.Div(
		h.Class("absolute left-full top-0 flex w-16 justify-center pt-5"),
		h.Button(
			h.Type("button"),
			h.OnClick(
				js.EvalCommandsOnSelector("#mobile-sidebar", js.RemoveClass("relative"), js.AddClass("hidden")),
			),
			h.Class("-m-2.5 p-2.5"),
			h.Span(
				h.Class("sr-only"),
				h.Text("Close sidebar"),
			),
			h.Svg(
				h.Class("size-6 text-white"),
				h.Attribute("fill", "none"),
				h.Attribute("viewBox", "0 0 24 24"),
				h.Attribute("stroke-width", "1.5"),
				h.Attribute("stroke", "currentColor"),
				h.Attribute("aria-hidden", "true"),
				h.Attribute("data-slot", "icon"),
				h.Path(
					h.Attribute("stroke-linecap", "round"),
					h.Attribute("stroke-linejoin", "round"),
					h.Attribute("d", "M6 18 18 6M6 6l12 12"),
				),
			),
		),
	)

	return h.Div(
		h.Class("hidden md:relative z-40 w-60"),
		h.Id("mobile-sidebar"),
		h.Role("dialog"),
		h.Attribute("aria-modal", "true"),
		h.Div(
			h.Class("fixed inset-0 bg-gray-900/80"),
			h.Attribute("aria-hidden", "true"),
		),
		CloseButton,
		h.Div(
			h.Class("fixed inset-y-0 z-40 flex w-60 flex-col"),
			h.Div(
				h.Class("flex grow flex-col gap-y-5 overflow-y-auto border-r border-gray-200 bg-white px-4 pb-4"),
				h.Div(
					h.Class("mt-2 flex h-10 shrink-0 items-center"),
					h.Div(
						h.Class("h-6 w-auto"),
						HtmgoLogo(),
					),
				),
				RoutingSection(),
				ResourceList(ctx),
				ServerList(ctx),
				DebugSection(),
			),
		),
	)
}

func ResourceList(ctx *h.RequestContext) *h.Element {
	list, err := app.ResourceList(ctx.ServiceLocator())

	if err != nil {
		list = []*app.Resource{}
	}

	return h.Div(
		h.Class("flex flex-col gap-2"),
		h.Div(
			h.Class("flex justify-between items-center"),
			h.P(
				h.Text("Resources"),
				h.Class("text-slate-800 font-bold"),
			),
			PrimaryButton(ButtonProps{
				Size: "xs",
				Text: "+ New",
				Href: urls.NewResourceUrl(),
			}),
			//h.A(
			//	h.Text("+ New"),
			//	h.Href(urls.NewResourceUrl()),
			//	h.Class("bg-slate-900 hover:bg-slate-800 text-white text-xs font-bold py-2 px-2 rounded"),
			//),
		),
		h.Div(
			h.Class("flex flex-col gap-2"),
			h.List(list, func(resource *app.Resource, index int) *h.Element {
				return h.A(
					h.Href(urls.ResourceUrl(resource.Id)),
					h.Text(resource.Name),
					h.Class("text-slate-900 hover:text-brand-400"),
				)
			}),
		),
	)
}

func ServerList(ctx *h.RequestContext) *h.Element {
	list, err := app.ServerList(ctx.ServiceLocator())

	if err != nil {
		list = []*app.Server{}
	}

	return h.Div(
		h.Class("flex flex-col gap-2"),
		h.Div(
			h.Class("flex justify-between items-center"),
			h.P(
				h.Text("Servers"),
				h.Class("text-slate-800 font-bold"),
			),
		),
		h.Div(
			h.Class("flex flex-col gap-2"),
			h.List(list, func(server *app.Server, index int) *h.Element {
				return h.A(
					h.Href(urls.ServerUrl(server.Id)),
					h.Text(h.Ternary(server.Name != "", server.Name, server.HostName)),
					h.Class("text-slate-900 hover:text-brand-400"),
				)
			}),
		),
	)
}

func RoutingSection() *h.Element {

	links := []Page{
		{
			Title: "Route Table",
			Path:  "/routing",
		},
	}

	return h.Div(
		h.Class("flex flex-col gap-2"),
		h.Div(
			h.Class("flex justify-between items-center"),
			h.P(
				h.Text("Routing"),
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
		{
			Title: "Interval Job Debug",
			Path:  "/debug/job",
		},
		{
			Title: "Router",
			Path:  "/debug/router",
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
