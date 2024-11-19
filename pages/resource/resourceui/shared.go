package resourceui

import (
	"errors"
	"github.com/maddalax/htmgo/framework/h"
	"paas/app"
	"paas/app/ui"
	"paas/app/urls"
	"paas/pages"
)

func Page(ctx *h.RequestContext, children func(resource *app.Resource) *h.Element) *h.Page {
	id := ctx.QueryParam("id")
	resource, err := app.ResourceGet(ctx.ServiceLocator(), id)

	if err != nil {
		ctx.Redirect("/", 302)
		return h.EmptyPage()
	}

	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-2 px-8"),
			PageHeader(ctx, resource),
			h.Div(
				h.GetPartialWithQs(GetStatusPartial, h.NewQs("id", resource.Id), "load, every 3s"),
			),
			children(resource),
		),
	)
}

func PageHeader(ctx *h.RequestContext, resource *app.Resource) *h.Element {
	_, err := app.IsResourceRunnable(resource)

	return h.Div(
		h.Class("flex flex-col gap-6"),
		h.If(err != nil, ResourceStatusError(err)),
		h.Id("resource-page-header"),
		h.Div(
			h.Class("flex gap-2 items-center"),
			h.H3F("%s", resource.Name, h.Class("text-2xl")),
			h.Div(
				h.Class("ml-0.5 mt-1"),
				ui.StatusIndicator(ui.StatusIndicatorProps{
					RunStatus: resource.RunStatus,
				}),
			),
		),
		TopTabs(ctx, resource, ui.LinkTabsProps{
			End: ResourceStatusContainer(resource),
		}),
	)
}

func ResourceStatusError(err error) *h.Element {
	if err == nil {
		return h.Empty()
	}
	switch {
	case errors.Is(err, app.DockerConnectionError):
		return ui.ErrorAlert(h.Pf("Failed to connect to docker"), h.Pf("Please check your docker connection"))
	}
	return ui.ErrorAlert(h.Pf("Failed to load resource status"), h.Pf(err.Error()))
}

func ResourceStatusContainer(resource *app.Resource) *h.Element {
	return h.Div(
		h.Id("resource-status"),
		h.Class("flex gap-2 absolute -top-3 right-0"),
		ResourceStatus(resource),
	)
}

func ResourceStatus(resource *app.Resource) *h.Element {
	runnable, err := app.IsResourceRunnable(resource)

	if err != nil {
		return h.Empty()
	}

	var deployButton = ui.SecondaryButton(ui.ButtonProps{
		Href: urls.ResourceStartDeploymentPath(resource.Id, ""),
		Text: "Deploy Resource",
	})

	var stopButton = ui.SubmitButton(ui.SubmitButtonProps{
		Post:           h.GetPartialPathWithQs(StopResource, h.NewQs("id", resource.Id)),
		SubmittingText: "Stopping...",
		Text:           "Stop",
	})

	var redeployButton = ui.PrimaryButton(ui.ButtonProps{
		Href: urls.ResourceStartDeploymentPath(resource.Id, ""),
		Text: "Redeploy",
	})

	var startButton = ui.SubmitButton(ui.SubmitButtonProps{
		Post:           h.GetPartialPathWithQs(StartResource, h.NewQs("id", resource.Id)),
		SubmittingText: "Starting...",
		Text:           "Start",
	})

	var restartButton = ui.SubmitButton(ui.SubmitButtonProps{
		Post:           h.GetPartialPathWithQs(RestartResource, h.NewQs("id", resource.Id)),
		SubmittingText: "Restarting...",
		Text:           "Restart",
	})

	return h.Div(
		h.Class("flex gap-2 w-full"),
		h.IfElse(!runnable, deployButton, redeployButton),
		h.If(resource.RunStatus == app.RunStatusRunning, stopButton),
		h.If(resource.RunStatus == app.RunStatusRunning, restartButton),
		h.If(resource.RunStatus != app.RunStatusRunning && runnable, startButton),
	)
}

func TopTabs(ctx *h.RequestContext, resource *app.Resource, props ui.LinkTabsProps) *h.Element {
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
