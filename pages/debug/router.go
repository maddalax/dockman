package debug

import (
	"dockside/app"
	"dockside/pages"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"strconv"
)

func RouterDebug(ctx *h.RequestContext) *h.Page {
	return pages.SidebarPage(ctx,
		h.Div(
			h.Class("p-4"),
			h.H3F(
				"Router Debug",
				h.Class("text-xl font-bold mb-6"),
			),
			h.Div(
				h.Class("overflow-x-auto"), // Ensure responsiveness for smaller screens
				h.Table(
					h.Class("w-full border-collapse border border-gray-300"),
					routerTableHeader(),
					h.TBody(
						h.GetPartial(RouterPartial, "load, every 3s"),
					),
				),
			),
		),
	)
}

func RouterPartial(ctx *h.RequestContext) *h.Partial {
	proxy := app.GetServiceRegistry(ctx.ServiceLocator()).GetReverseProxy()
	upstreams := proxy.GetUpstreams()

	return h.NewPartial(
		h.List(upstreams, func(item *app.CustomUpstream, index int) *h.Element {
			return routerTableRow(item)
		}),
	)
}

func routerTableHeader() *h.Element {
	return h.THead(
		h.Tr(
			h.Class("bg-gray-100 text-left border-b border-gray-300"),
			tableHeaderCell("Server"),
			tableHeaderCell("Upstreams"),
			tableCell("Resource"),
			tableHeaderCell("Total Requests"),
			tableHeaderCell("Last Request"),
			tableHeaderCell("Avg Response Time"),
			tableHeaderCell("Status"),
			tableHeaderCell("Match"),
		),
	)
}

func routerTableRow(metric *app.CustomUpstream) *h.Element {
	return h.Tr(
		h.Class("border-b border-gray-300 hover:bg-gray-50"),
		tableCell(metric.Metadata.Server.FormattedName()),
		tableCell(metric.Url.String()),
		tableCell(metric.Metadata.Resource.Name),
		tableCell(strconv.Itoa(int(metric.TotalRequests.Load()))),
		tableCell(formatTimePretty(metric.LastRequest)),
		tableCell(fmt.Sprintf("%dms", metric.AverageResponseTime.Milliseconds())),
		tableCell(h.Ternary(metric.Healthy, "Healthy", "Unhealthy")),
		h.Td(
			h.Class("py-2 px-4 text-sm text-gray-700"),
			upstreamBlockView(metric),
		),
	)
}

func upstreamBlockView(metric *app.CustomUpstream) *h.Element {
	block := metric.Metadata.Block
	return h.Div(
		h.Pf("Host: %s", block.Hostname),
		h.Pf("Path: %s", block.Path),
		h.Pf("Mod: %s", block.PathMatchModifier),
	)
}
