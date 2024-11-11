package ui

import (
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
)

func LogBody() *h.Element {
	return h.Div(
		h.Class("max-w-[800px] max-h-full overflow-y-auto bg-white border border-gray-300 rounded-lg shadow-lg p-4 mt-6 min-w-[800px]"),
		h.Id("build-log"),
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
	return h.Div(
		h.Class("bg-slate-50 p-2 rounded-md text-sm"),
		h.Attribute("hx-swap-oob", "beforeend:#build-log"),
		h.P(
			h.Text(data),
		),
	)
}
