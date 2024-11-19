package ui

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/internal"
)

type StatusIndicatorProps struct {
	RunStatus internal.RunStatus
}

func StatusIndicator(props StatusIndicatorProps) h.Ren {
	var colorClass string
	var animationClass string

	if props.RunStatus == internal.RunStatusRunning {
		colorClass = "bg-green-500"
		animationClass = "animate-pulse"
	} else if props.RunStatus == internal.RunStatusPartiallyRunning {
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
			h.TextF("%s", statusText(props.RunStatus)),
		),
	)
}

// statusText returns the textual representation of the status.
func statusText(status internal.RunStatus) string {
	switch status {
	case internal.RunStatusRunning:
		return "Running"
	case internal.RunStatusPartiallyRunning:
		return "Partially Running"
	default:
		return "Stopped"
	}
}
