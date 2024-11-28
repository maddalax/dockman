package debug

import (
	"dockman/app"
	"dockman/app/ui"
	"dockman/pages"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"strconv"
)

func RouterDebug(ctx *h.RequestContext) *h.Page {
	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("p-4"),
			h.H3F(
				"Router Debug",
				h.Class("text-xl font-bold mb-6"),
			),
			h.Div(
				h.Class("overflow-x-auto"),
				// Ensure responsiveness for smaller screens
				h.Div(
					h.GetPartial(RouterPartial, "load, every 3s"),
				),
			),
		),
	)
}

func RouterPartial(ctx *h.RequestContext) *h.Partial {
	proxy := app.GetServiceRegistry(ctx.ServiceLocator()).GetReverseProxy()
	upstreams := proxy.GetUpstreams()

	table := ui.NewTable()

	table.AddColumns([]string{
		"Server",
		"Upstreams",
		"Resource",
		"Total Requests",
		"Last Request",
		"Avg Response Time",
		"Status",
		"Match",
	})

	for _, upstream := range upstreams {
		table.AddRow()

		table.WithCellTexts(
			upstream.Metadata.Server.FormattedName(),
			upstream.Url.String(),
			upstream.Metadata.Resource.Name,
			strconv.Itoa(int(upstream.TotalRequests.Load())),
			formatTimePretty(upstream.LastRequest),
			fmt.Sprintf("%dms", upstream.AverageResponseTime.Milliseconds()),
			h.Ternary(upstream.Healthy, "Healthy", "Unhealthy"),
		)

		table.AddCell(upstreamBlockView(upstream))
	}

	return h.NewPartial(
		table.Render(),
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
