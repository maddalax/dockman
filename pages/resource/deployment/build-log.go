package deployment

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/pages"
	resource2 "paas/pages/resource"
	"paas/resources"
	"paas/ui"
)

func BuidLog(ctx *h.RequestContext) *h.Page {
	id := ctx.QueryParam("id")
	buildId := ctx.QueryParam("buildId")

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
			resource2.TopTabs(ctx, resource),
			h.Div(
				ui.DockerBuildLogs(ctx, resource, buildId),
			),
		),
	)
}
