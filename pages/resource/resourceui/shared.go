package resourceui

import (
	"dockside/app"
	"dockside/app/ui"
	"dockside/app/urls"
	"dockside/pages"
	"errors"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
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
	locator := ctx.ServiceLocator()
	_, err := app.IsResourceRunnable(locator, resource)

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
					RunStatus: app.GetComputedRunStatus(resource),
				}),
			),
		),
		TopTabs(ctx, resource, ui.LinkTabsProps{
			End: ResourceStatusContainer(locator, resource),
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

func ResourceStatusContainer(locator *service.Locator, resource *app.Resource) *h.Element {
	return h.Div(
		h.Id("resource-status"),
		h.Class("flex gap-2 absolute -top-3 right-0"),
		ResourceStatus(locator, resource),
	)
}

func ResourceStatus(locator *service.Locator, resource *app.Resource) *h.Element {
	runnable, err := app.IsResourceRunnable(locator, resource)
	runStatus := app.GetComputedRunStatus(resource)

	if err != nil {
		return h.Empty()
	}

	var deployButton = ui.PrimaryButton(ui.ButtonProps{
		Href: urls.ResourceStartDeploymentPath(resource.Id, ""),
		Text: "Deploy Resource",
	})

	var stopButton = ui.DangerButton(ui.ButtonProps{
		Post:           h.GetPartialPathWithQs(StopResource, h.NewQs("id", resource.Id)),
		SubmittingText: "Stopping...",
		Text:           "Stop",
	})

	var redeployButton = ui.SecondaryButton(ui.ButtonProps{
		Href: urls.ResourceStartDeploymentPath(resource.Id, ""),
		Text: "Redeploy",
	})

	var startButton = ui.SubmitButton(ui.ButtonProps{
		Post:           h.GetPartialPathWithQs(StartResource, h.NewQs("id", resource.Id)),
		SubmittingText: "Starting...",
		Text:           "Start",
	})

	var restartButton = ui.SubmitButton(ui.ButtonProps{
		Post:           h.GetPartialPathWithQs(RestartResource, h.NewQs("id", resource.Id)),
		SubmittingText: "Restarting...",
		Text:           "Restart",
	})

	return h.Div(
		h.Class("flex gap-2 w-full"),
		h.IfElse(!runnable, deployButton, redeployButton),
		h.If(runStatus != app.RunStatusNotRunning, stopButton),
		h.If(runStatus == app.RunStatusRunning || runStatus == app.RunStatusPartiallyRunning, restartButton),
		h.If(runStatus == app.RunStatusNotRunning && runnable, startButton),
	)
}

func TopTabs(ctx *h.RequestContext, resource *app.Resource, props ui.LinkTabsProps) *h.Element {
	props.Links = []ui.Link{
		{
			Text: "Overview",
			Href: urls.ResourceUrl(resource.Id),
		},
		{
			Text: "Servers",
			Href: urls.ResourceServersUrl(resource.Id),
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
