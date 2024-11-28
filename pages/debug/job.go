package debug

import (
	"dockside/app"
	"dockside/pages"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"slices"
	"strconv"
	"strings"
	"time"
)

func IntervalJobDebug(ctx *h.RequestContext) *h.Page {
	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("p-4"),
			h.Div(
				h.Class("flex flex-col gap-1 mb-4"),
				h.H3F(
					"Interval Job Debug",
					h.Class("text-xl font-bold"),
				),
				h.Div(
					h.GetPartial(LastUpdatedPartial, "load, every 3s"),
				),
			),
			h.Div(
				h.Class("overflow-x-auto"),
				// Ensure responsiveness for smaller screens
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

func LastUpdatedPartial(ctx *h.RequestContext) *h.Partial {
	manager := app.GetServiceRegistry(ctx.ServiceLocator()).GetJobMetricsManager()
	metrics := manager.GetMetrics()

	lastUpdated := time.Time{}

	for _, metric := range metrics {
		if metric.LastRan.After(lastUpdated) {
			lastUpdated = metric.LastRan
		}
	}

	return h.NewPartial(
		h.Div(
			h.Class("text-sm text-gray-500"),
			h.Text("Last updated: "),
			h.Text(lastUpdated.Format("Jan 2, 2006 at 3:04:05 PM")),
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
		keyA := strings.Join([]string{a.JobSource, a.JobName}, "-")
		keyB := strings.Join([]string{b.JobSource, b.JobName}, "-")
		return strings.Compare(keyA, keyB)
	})

	return h.NewPartial(
		h.List(metrics, func(item *app.JobMetric, index int) *h.Element {
			return tableRow(item)
		}),
	)
}

func ToggleJob(ctx *h.RequestContext) *h.Partial {
	jobName := ctx.QueryParam("job")
	if jobName == "" {
		return h.EmptyPartial()
	}
	runner := app.GetServiceRegistry(ctx.ServiceLocator()).GetJobRunner()
	job := runner.GetJob(jobName)
	if job == nil {
		return h.EmptyPartial()
	}
	job.Toggle()
	return h.EmptyPartial()
}

func tableHeader() *h.Element {
	return h.THead(
		h.Tr(
			h.Class("bg-gray-100 text-left border-b border-gray-300"),
			tableHeaderCell("Source"),
			tableHeaderCell("Status"),
			tableHeaderCell("Job Name"),
			tableHeaderCell("Description"),
			tableHeaderCell("Interval"),
			tableHeaderCell("Last Ran"),
			tableHeaderCell("Total Runs"),
			tableHeaderCell("Last Run Duration"),
			tableHeaderCell("Actions"),
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
		tableCell(metric.JobSource),
		tableCell(
			metric.Status,
			h.Ternary(metric.Status == "running", "text-green-700", "text-red-700"),
		),
		tableCell(metric.JobName),
		tableCell(metric.JobDescription),
		tableCell(metric.Interval.String()),
		tableCell(formatTimePretty(metric.LastRan)),
		tableCell(strconv.Itoa(metric.TotalRuns)),
		tableCell(metric.LastRunDuration.String()),
		h.Td(
			h.Class("py-2 px-4 text-sm text-gray-700"),
			// can only pause or resume dockside jobs, not server jobs
			h.If(
				metric.JobSource == "dockside",
				h.Button(
					h.NoSwap(),
					h.PostPartialWithQs(
						ToggleJob,
						h.NewQs("job", fmt.Sprintf("%s-%s", metric.JobSource, metric.JobName)),
					),
					h.Text(
						h.Ternary(metric.JobPaused, "Resume", "Pause"),
					),
					h.Class("text-blue-500 hover:text-blue-700"),
				),
			),
		),
	)
}

func tableCell(value string, classes ...string) *h.Element {
	classes = append(classes, "py-2 px-4 text-sm text-gray-700")
	return h.Td(
		h.Class(classes...),
		h.Pf(
			value,
			h.Class("truncate"),
		),
	)
}

func calculateRunStatus(metric *app.JobMetric) string {
	if metric.JobPaused {
		return "paused"
	}

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
