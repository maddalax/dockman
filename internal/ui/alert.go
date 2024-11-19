package ui

import (
	"github.com/maddalax/htmgo/framework/h"
)

func ErrorAlertPartial(ctx *h.RequestContext, title *h.Element, message *h.Element) *h.Partial {
	return h.SwapPartial(
		ctx,
		ErrorAlert(title, message),
	)
}

func GenericErrorAlertPartial(ctx *h.RequestContext, err error) *h.Partial {
	return ErrorAlertPartial(
		ctx,
		h.Pf("Unable to perform the operation"),
		h.Pf(err.Error()),
	)
}

func SuccessAlertPartial(ctx *h.RequestContext, title string, message string) *h.Partial {
	return h.SwapPartial(
		ctx,
		SuccessAlert(
			h.Pf(title),
			h.Pf(message),
		),
	)
}

func AlertPlaceholder() *h.Element {
	return h.Div(
		h.Id("ui-alert"),
	)
}

func SuccessAlert(title *h.Element, message *h.Element) *h.Element {
	return h.Div(
		h.Id("ui-alert"),
		h.Role("alert"),
		h.Class("rounded border-s-4 border-green-500 bg-green-50 p-4 w-full"),
		h.Strong(
			h.Class("block font-medium text-green-800"),
			title,
		),
		h.P(
			h.Class("mt-2 text-sm text-green-700"),
			message,
		),
	)
}

func ErrorAlert(title *h.Element, message *h.Element) *h.Element {
	return h.Div(
		h.Id("ui-alert"),
		h.Role("alert"),
		h.Class("rounded border-s-4 border-red-500 bg-red-50 p-4 w-full"),
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
