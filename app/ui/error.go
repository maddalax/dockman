package ui

import "github.com/maddalax/htmgo/framework/h"

func FormError(error string) *h.Element {
	return h.Div(
		h.Id("form-error"),
		h.If(
			error != "",
			ErrorAlert(h.Pf(error), nil),
		),
	)
}

func SwapFormError(ctx *h.RequestContext, error string) *h.Partial {
	return h.SwapPartial(ctx,
		FormError(error),
	)
}
