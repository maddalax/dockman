package ui

import (
	"dockside/app"
	"dockside/app/util"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
	"strings"
)

type LogBodyOptions struct {
	MaxLogs int
}

func LogBody(opts LogBodyOptions) *h.Element {
	return h.Div(
		h.Class("w-full max-h-full h-full overflow-y-auto bg-white border border-gray-300 rounded-lg shadow-lg mt-6 bg-red-500"),
		h.Div(
			h.Id("build-log"),
			h.Class("flex flex-col w-full"),
		),
		// Scroll to the bottom when the page loads
		h.OnLoad(
			js.EvalJs(`
				setTimeout(() => {
					const logs = document.getElementById('build-log');
					logs.parentElement.scrollTop = logs.parentElement.scrollHeight;
				}, 1000)
			`),
		),
		// Scroll to the bottom when a new message is added
		h.OnEvent("htmx:wsAfterMessage",
			js.EvalJs(fmt.Sprintf(`
				// Remove excess logs
				const logs = document.getElementById('build-log');
				while (logs.children.length >= %d) {
					logs.removeChild(logs.firstElementChild);
				}
			`, opts.MaxLogs)),
			js.EvalJs(`
				const logsContainer = document.getElementById('build-log').parentElement;
				const scrollPosition = logsContainer.scrollTop + logsContainer.clientHeight;
				const distanceFromBottom = logsContainer.scrollHeight - scrollPosition;
				const scrollThreshold = 1000; // Adjust this as needed
				if (distanceFromBottom <= scrollThreshold) {
					logsContainer.scrollTop = logsContainer.scrollHeight;
				}
			`),
		),
	)
}

func DockerLogLine(log *app.DockerLog) *h.Element {
	swap := h.Attribute("hx-swap-oob", "beforeend:#build-log")

	return h.Div(
		swap,
		h.Div(
			h.Class("px-4 flex flex-no-wrap items-start gap-4 border-b border-gray-300 py-1"),
			h.Div(
				h.Class("w-1/8 truncate text-sm font-medium"),
				h.Text(log.HostName),
			),
			h.Div(
				h.Class("w-1/8 text-sm text-gray-600"),
				h.Text(log.Time.Format("2006-01-02 15:04:05")),
			),
			h.Div(
				h.Class("flex-1 text-sm text-gray-800"),
				h.Text(log.Log),
			),
		),
	)
}

func LogLine(data string) *h.Element {
	swap := h.Attribute("hx-swap-oob", "beforeend:#build-log")

	if strings.HasPrefix(data, "BUILD_ERROR:") {
		data = strings.TrimPrefix(data, "BUILD_ERROR:")
		return h.Div(
			swap,
			h.Div(
				h.Class("px-4 flex items-start gap-4 border-b border-red-300 py-1"),
				h.Div(
					h.Class("w-1/8 truncate text-sm font-medium text-red-600"),
					h.Text("Error"),
				),
				h.Div(
					h.Class("flex-1 text-sm text-red-800"),
					h.UnsafeRaw(util.Sanitize(data)),
				),
			),
		)
	}

	return h.Div(
		swap,
		h.Div(
			h.Class("px-4 flex items-start gap-4 border-b border-gray-300 py-1"),
			h.Div(
				h.Class("flex-1 text-sm text-gray-800"),
				h.UnsafeRaw(util.Sanitize(data)),
			),
		),
	)
}
