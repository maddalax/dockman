package ui

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
	"paas/util"
	"strings"
)

type LogBodyOptions struct {
	MaxLogs int
}

func LogBody(opts LogBodyOptions) *h.Element {
	return h.Div(
		h.Class("max-w-4xl max-h-full overflow-y-auto bg-white border border-gray-300 rounded-lg shadow-lg p-4 mt-6 min-w-[800px]"),
		h.Div(
			h.Id("build-log"),
		),
		// Scroll to the bottom of the div when the page loads
		h.OnLoad(
			// language=JavaScript
			js.EvalJs(`
					setTimeout(() => {
           self.scrollTop = self.scrollHeight;       
					}, 1000)
				`),
		),
		// Scroll to the bottom of the div when the message is sent
		// only if the user is close to the bottom of the div
		h.OnEvent("htmx:wsAfterMessage",
			// language=JavaScript
			js.EvalJs(fmt.Sprintf(
				`
				 // only keep the last MaxLogs
         let logs = document.getElementById('build-log');
				 while (logs.children.length >= %d) {
       		 logs.removeChild(logs.firstElementChild);
    			}
      `,
				opts.MaxLogs)),
			// language=JavaScript
			js.EvalJs(`
					const scrollPosition = self.scrollTop + self.clientHeight;
    			const distanceFromBottom = self.scrollHeight - scrollPosition;
    			const scrollThreshold = 1000; // Adjust this to define how close the user should be to the bottom
					 if (distanceFromBottom <= scrollThreshold) {
        			self.scrollTop = self.scrollHeight;
    			}
				`),
		),
	)
}

func LogLine(data string) *h.Element {
	swap := h.Attribute("hx-swap-oob", "beforeend:#build-log")

	if strings.HasPrefix(data, "BUILD_ERROR:") {
		data = strings.TrimPrefix(data, "BUILD_ERROR:")
		return h.Div(
			swap,
			h.P(
				h.Class("text-red-800"),
				h.UnsafeRaw(util.Sanitize(data)),
			),
		)
	}

	return h.Div(
		swap,
		h.P(
			h.UnsafeRaw(util.Sanitize(data)),
		),
	)
}
