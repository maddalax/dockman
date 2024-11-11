package resource

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/pages"
	"paas/resources"
)

func Index(ctx *h.RequestContext) *h.Page {
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
			h.Pf("Environment: %s", resource.Environment),
			h.Pf("Run Type: %v", resource.RunType),
			h.Pf("Build Meta: %v", resource.BuildMeta),
			h.Pf("Env: %s", resource.Env),
		),
	)
}
