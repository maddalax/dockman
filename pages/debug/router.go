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
	config := proxy.GetConfig()
	return h.NewPartial(
		h.List(config.Upstreams, func(item *app.UpstreamWithResource, index int) *h.Element {
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
		),
	)
}

func routerTableRow(metric *app.UpstreamWithResource) *h.Element {
	return h.Tr(
		h.Class("border-b border-gray-300 hover:bg-gray-50"),
		tableCell(metric.Server.FormattedName()),
		tableCell(metric.Upstream.Url.String()),
		tableCell(metric.Resource.Name),
		tableCell(strconv.Itoa(metric.Upstream.TotalRequests)),
		tableCell(formatTimePretty(metric.Upstream.LastRequest)),
		tableCell(fmt.Sprintf("%dms", metric.Upstream.AverageResponseTime.Milliseconds())),
		tableCell(h.Ternary(metric.Upstream.Healthy, "Healthy", "Unhealthy")),
	)
}
