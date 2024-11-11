package resource

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/pages"
)

func Index(ctx *h.RequestContext) *h.Page {
	id := ctx.QueryParam("id")
	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-2"),
			pages.Title("Resource"),
			h.Pf("Resource ID: %s", id),
		),
	)
}
