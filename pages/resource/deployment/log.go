package deployment

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/pages"
	"paas/resources"
	"paas/ui"
	"paas/urls"
)

func NewDeployment(ctx *h.RequestContext) *h.Page {
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
			ui.LinkTabs(ctx, ui.LinkTabsProps{
				Links: []ui.Link{
					{
						Text: "Overview",
						Href: urls.ResourceUrl(resource.Id),
					},
					{
						Text: "Deployment",
						Href: urls.ResourceDeploymentUrl(resource.Id),
					},
					{
						Text: "Environment",
						Href: urls.ResourceEnvironmentUrl(resource.Id),
					},
				},
			}),
			ui.DockerBuildLogs(ctx, resource, buildId),
		),
	)
}
