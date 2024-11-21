package debug

import (
	"dockside/app"
	"dockside/pages"
	"github.com/maddalax/htmgo/framework/h"
	"strconv"
	"time"
)

func IntervalJobDebug(ctx *h.RequestContext) *h.Page {
	return pages.SidebarPage(ctx,
		h.Div(
			h.Class("p-4"),
			h.H3F(
				"Interval Job Debug",
				h.Class("text-xl font-bold mb-6"),
			),
			h.Div(
				h.Class("overflow-x-auto"), // Ensure responsiveness for smaller screens
				h.Table(
					h.Class("w-full border-collapse border border-gray-300"),
					tableHeader(),
					h.TBody(
						h.GetPartial(JobMetricsPartial, "load, every 3s"),
					),
				),
			),
		),
	)
}

func JobMetricsPartial(ctx *h.RequestContext) *h.Partial {
	manager := app.NewJobMetricsManager(ctx.ServiceLocator())
	metrics := manager.GetMetrics()

	return h.NewPartial(
		h.List(metrics, func(item *app.JobMetric, index int) *h.Element {
			return tableRow(item)
		}),
	)
}

func tableHeader() *h.Element {
	return h.THead(
		h.Tr(
			h.Class("bg-gray-100 text-left border-b border-gray-300"),
			tableHeaderCell("Job Name"),
			tableHeaderCell("Status"),
			tableHeaderCell("Last Ran"),
			tableHeaderCell("Total Runs"),
			tableHeaderCell("Last Run Duration"),
		),
	)
}

func tableHeaderCell(label string) *h.Element {
	return h.Th(
		h.Class("py-2 px-4 text-sm font-semibold text-gray-700"),
		h.Text(label),
	)
}

func tableRow(metric *app.JobMetric) *h.Element {
	return h.Tr(
		h.Class("border-b border-gray-300 hover:bg-gray-50"),
		tableCell(metric.JobName),
		tableCell(metric.Status),
		tableCell(formatTimePretty(metric.LastRan)),
		tableCell(strconv.Itoa(metric.TotalRuns)),
		tableCell(metric.LastRunDuration.String()),
	)
}

func tableCell(value string) *h.Element {
	return h.Td(
		h.Class("py-2 px-4 text-sm text-gray-700"),
		h.Text(value),
	)
}

func formatTimePretty(t time.Time) string {
	return t.Format("Jan 2, 2006 at 3:04 PM")
}
