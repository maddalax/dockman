package pages

import (
	"github.com/maddalax/htmgo/framework/h"
)

func Index(ctx *h.RequestContext) *h.Page {
	return SidebarPage(
		ctx,
		h.Div(),
	)
}
