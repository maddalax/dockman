package resourceui

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/domain"
	"paas/pages"
	"paas/resources"
	"paas/ui"
	"paas/urls"
)

func Page(ctx *h.RequestContext, children func(resource *domain.Resource) *h.Element) *h.Page {
	id := ctx.QueryParam("id")
	resource, err := resources.Get(ctx.ServiceLocator(), id)

	if err != nil {
		ctx.Redirect("/", 302)
		return h.EmptyPage()
	}

	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-2 px-8"),
			PageHeader(ctx, resource),
			children(resource),
		),
	)
}

func PageHeader(ctx *h.RequestContext, resource *domain.Resource) *h.Element {
	return h.Div(
		h.Class("flex flex-col gap-6"),
		h.Div(
			h.Class("flex gap-2 justify-between items-center"),
			h.H3F("%s", resource.Name, h.Class("text-2xl")),
		),
		TopTabs(ctx, resource, ui.LinkTabsProps{
			End: h.Div(
				h.GetPartialWithQs(GetStatusPartial, h.NewQs("id", resource.Id), "load"),
			),
		}),
	)
}

func ResourceStatusContainer(resource *domain.Resource) *h.Element {
	return h.Div(
		h.Id("resource-status"),
		h.Class("flex gap-2 absolute -top-3 right-0"),
		ResourceStatus(resource),
	)
}

func ResourceStatus(resource *domain.Resource) *h.Element {
	if resource.RunStatus == domain.RunStatusRunning {
		return ui.SubmitButton(ui.SubmitButtonProps{
			Post:           h.GetPartialPathWithQs(StopResource, h.NewQs("id", resource.Id)),
			SubmittingText: "Stopping...",
			Text:           "Stop Resource",
		})
	}

	return ui.SubmitButton(ui.SubmitButtonProps{
		Text:           "Start Resource",
		SubmittingText: "Starting...",
		Post:           h.GetPartialPathWithQs(StartResource, h.NewQs("id", resource.Id)),
	})
}

func TopTabs(ctx *h.RequestContext, resource *domain.Resource, props ui.LinkTabsProps) *h.Element {
	props.Links = []ui.Link{
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
	}

	return ui.LinkTabs(ctx, props)
}
