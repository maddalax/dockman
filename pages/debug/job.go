package debug

import (
	"dockside/app"
	"dockside/pages"
	"github.com/maddalax/htmgo/framework/h"
	"slices"
	"strconv"
	"strings"
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
	manager := app.GetServiceRegistry(ctx.ServiceLocator()).GetJobMetricsManager()
	metrics := manager.GetMetrics()

	for i, metric := range metrics {
		metrics[i].Status = calculateRunStatus(metric)
	}

	slices.SortFunc(metrics, func(a, b *app.JobMetric) int {
		// Compare by Status first
		if cmp := strings.Compare(a.Status, b.Status); cmp != 0 {
			return cmp
		}
		// If Status is equal, compare by Name
		return strings.Compare(a.JobName, b.JobName)
	})

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
			tableHeaderCell("Interval"),
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
		tableCell(metric.Interval.String()),
		tableCell(metric.Status, h.Ternary(metric.Status == "running", "text-green-700", "text-red-700")),
		tableCell(formatTimePretty(metric.LastRan)),
		tableCell(strconv.Itoa(metric.TotalRuns)),
		tableCell(metric.LastRunDuration.String()),
	)
}

func tableCell(value string, classes ...string) *h.Element {
	classes = append(classes, "py-2 px-4 text-sm text-gray-700")
	return h.Td(
		h.Class(classes...),
		h.Text(value),
	)
}

func calculateRunStatus(metric *app.JobMetric) string {
	timeBetween := metric.Interval - metric.LastRunDuration
	buffer := timeBetween / 10 // Add a 10% buffer
	adjustedTime := timeBetween + buffer
	hasNotRunInTime := metric.LastRan.Before(time.Now().Add(-adjustedTime))

	if hasNotRunInTime {
		return "stopped"
	}
	return "running"
}

func formatTimePretty(t time.Time) string {
	return t.Format("Jan 2, 2006 at 3:04 PM")
}
