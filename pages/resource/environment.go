package resource

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/pages"
	"paas/resources"
)

func Environment(ctx *h.RequestContext) *h.Page {
	id := ctx.QueryParam("id")

	resource, err := resources.Get(ctx.ServiceLocator(), id)

	if err != nil {
		ctx.Redirect("/", 302)
		return h.EmptyPage()
	}

	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-2"),
			pages.Title("Resource"),
			h.Pf("Resource: %s", resource.Name),
			TopTabs(ctx, resource),
		),
	)
}
