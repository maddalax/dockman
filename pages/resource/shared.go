package resource

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/resources"
	"paas/ui"
	"paas/urls"
)

func TopTabs(ctx *h.RequestContext, resource *resources.Resource) *h.Element {
	return ui.LinkTabs(ctx, ui.LinkTabsProps{
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
			{
				Text: "Run Log",
				Href: urls.ResourceRunLogUrl(resource.Id),
			},
		},
	})
}
