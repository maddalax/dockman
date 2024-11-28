package debug

import (
	"dockman/app"
	"dockman/app/ui"
	"dockman/pages"
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
				h.GetPartial(JobMetricsPartial, "load, every 3s"),
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

	table := ui.NewTable()

	table.AddColumns([]string{
		"Source",
		"Status",
		"Job Name",
		"Description",
		"Interval",
		"Last Ran",
		"Total Runs",
		"Last Run Duration",
		"Actions",
	})

	for _, metric := range metrics {
		table.AddRow()

		table.WithCellTexts(
			metric.JobSource,
			metric.Status,
			metric.JobName,
			metric.JobDescription,
			metric.Interval.String(),
			formatTimePretty(metric.LastRan),
			strconv.Itoa(metric.TotalRuns),
			metric.LastRunDuration.String(),
		)

		table.AddCell(
			h.Ternary(
				metric.JobSource == "dockman",
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
				h.Empty(),
			),
		)
	}

	return h.NewPartial(table.Render())
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
