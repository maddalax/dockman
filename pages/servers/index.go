package servers

import (
	"dockman/app"
	"dockman/app/ui"
	"dockman/pages"
	"github.com/maddalax/htmgo/framework/h"
)

func ServerListPartial(ctx *h.RequestContext) *h.Partial {
	return h.SwapPartial(ctx, serverList(ctx))
}

func ServerListPage(ctx *h.RequestContext) *h.Page {
	return pages.SidebarPage(
		ctx,
		h.Div(
			h.H1F("Servers"),
			h.Div(
				h.GetPartial(ServerListPartial, "load, every 3s"),
			),
			h.Div(
				h.Id("server-list"),
			),
		),
	)
}

func serverList(ctx *h.RequestContext) *h.Element {
	servers, err := app.ServerList(ctx.ServiceLocator())

	if err != nil {
		servers = []*app.Server{}
	}

	return h.Div(
		h.Id("server-list"),
		h.Class("flex flex-col gap-4"),
		h.List(servers, func(server *app.Server, index int) *h.Element {
			return ServerBlockDetails(server)
		}),
	)
}

func ServerBlockDetails(server *app.Server) *h.Element {
	runStatus := app.RunStatusNotRunning
	if server.IsAccessible() {
		runStatus = app.RunStatusRunning
	}

	return h.Div(
		h.Class("bg-white shadow-md rounded-lg p-4"),
		h.P(
			h.Span(
				h.Class("font-bold"),
				h.Text("Host Name: "),
			),
			h.Text(server.HostName),
			h.Class("text-slate-800"),
		),
		h.If(
			server.Name != "",
			h.P(
				h.Text(server.Name),
				h.Class("text-slate font-bold"),
			),
		),
		h.P(
			h.Span(
				h.Class("font-bold"),
				h.Text("Local IP Address: "),
			),
			h.Text(server.LocalIpAddress),
			h.Class("text-slate-800"),
		),
		h.P(
			h.Span(
				h.Class("font-bold"),
				h.Text("Remote IP Address: "),
			),
			h.Text(server.RemoteIpAddress),
			h.Class("text-slate-800"),
		),
		h.P(
			h.Span(
				h.Class("font-bold"),
				h.Text("OS: "),
			),
			h.Text(server.Os),
			h.Class("text-slate-800"),
		),
		h.P(
			h.Span(
				h.Class("font-bold"),
				h.Text("Last Seen: "),
			),
			h.Text(server.LastSeen.Format("2006-01-02 15:04:05")),
			h.Class("text-slate-800"),
		),
		h.P(
			ui.StatusIndicator(ui.StatusIndicatorProps{
				RunStatus: runStatus,
				TextMap: map[app.RunStatus]string{
					app.RunStatusNotRunning: "Not Accessible",
					app.RunStatusRunning:    "Connected",
				},
			}),
			h.Class("text-slate-800"),
		),
	)
}
