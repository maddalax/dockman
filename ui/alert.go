package ui

import (
	"github.com/maddalax/htmgo/framework/h"
)

func ErrorAlert(title *h.Element, message *h.Element) *h.Element {
	return h.Div(
		h.Role("alert"),
		h.Class("rounded border-s-4 border-red-500 bg-red-50 p-4"),
		h.Strong(
			h.Class("block font-medium text-red-800"),
			title,
		),
		h.P(
			h.Class("mt-2 text-sm text-red-700"),
			message,
		),
	)
}
