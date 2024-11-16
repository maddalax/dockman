package ui

import "github.com/maddalax/htmgo/framework/h"

type StatusIndicatorProps struct {
	IsRunning bool
}

func StatusIndicator(props StatusIndicatorProps) h.Ren {
	var colorClass string
	var animationClass string

	if props.IsRunning {
		colorClass = "bg-green-500"
		animationClass = "animate-pulse" // Tailwind animation class for slow pulsing
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
			h.TextF("%s", statusText(props.IsRunning)),
		),
	)
}

// statusText returns the textual representation of the status.
func statusText(isRunning bool) string {
	if isRunning {
		return "Running"
	}
	return "Stopped"
}
