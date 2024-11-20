package ui

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/app"
)

type StatusIndicatorProps struct {
	RunStatus app.RunStatus
	TextMap   map[app.RunStatus]string
}

func StatusIndicator(props StatusIndicatorProps) h.Ren {
	var colorClass string
	var animationClass string

	if props.RunStatus == app.RunStatusRunning {
		colorClass = "bg-green-500"
		animationClass = "animate-pulse"
	} else if props.RunStatus == app.RunStatusPartiallyRunning {
		colorClass = "bg-amber-500"
		animationClass = "animation-pulse"
	} else {
		colorClass = "bg-red-500"
		animationClass = "" // No animation for stopped
	}

	return h.Div(
		h.Class("flex items-center space-x-1"),
		h.Span(
			h.Class("h-3 w-3 rounded-full "+colorClass+" "+animationClass),
		),
		h.Span(
			h.Class("text-sm"),
			h.TextF("%s", h.Ternary(len(props.TextMap) > 0, props.TextMap[props.RunStatus], statusText(props.RunStatus))),
		))
}

// statusText returns the textual representation of the status.
func statusText(status app.RunStatus) string {
	switch status {
	case app.RunStatusRunning:
		return "Running"
	case app.RunStatusPartiallyRunning:
		return "Partially Running"
	default:
		return "Stopped"
	}
}
