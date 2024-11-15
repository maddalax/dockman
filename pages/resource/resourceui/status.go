package resourceui

import (
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"paas/docker"
	"paas/monitor"
	"paas/resources"
	"paas/ui"
)

func GetStatusPartial(ctx *h.RequestContext) *h.Partial {
	id := ctx.QueryParam("id")
	resource, err := resources.Get(ctx.ServiceLocator(), id)

	if err != nil {
		// TODO
		panic(err)
	}

	return h.NewPartial(
		ResourceStatusContainer(resource),
	)
}

func StartResource(ctx *h.RequestContext) *h.Partial {
	id := ctx.QueryParam("id")

	if id == "" {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf("Resource ID is required"), h.Empty()),
		))
	}

	client, err := docker.Connect()

	if err != nil {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf("Failed to connect to docker daemon"), h.Empty()),
		))
	}

	resource, err := resources.Get(ctx.ServiceLocator(), id)

	if err != nil {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf(err.Error()), h.Empty()),
		))
	}

	err = client.Run(resource, docker.RunOptions{
		KillExisting: true,
	})

	if err != nil {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf(err.Error()), h.Empty()),
		))
	}

	m := service.Get[monitor.Monitor](ctx.ServiceLocator())
	resource.RunStatus = m.GetRunStatus(resource)
	_ = resource.SetRunStatus(ctx.ServiceLocator(), resource.RunStatus)

	return h.SwapPartial(ctx, ResourceStatusContainer(resource))
}
