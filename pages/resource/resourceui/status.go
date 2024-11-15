package resourceui

import (
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"paas/domain"
	"paas/monitor"
	"paas/resources"
	"paas/ui"
	"paas/util"
	"time"
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
	resource, err := resources.Start(ctx.ServiceLocator(), id)

	if err != nil {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf(err.Error()), h.Empty()),
		))
	}

	// wait until its running
	m := service.Get[monitor.Monitor](ctx.ServiceLocator())

	success := util.WaitFor(time.Second*5, func() bool {
		return m.GetRunStatus(resource) == domain.RunStatusRunning
	})

	if !success {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf("Failed to start resource after 5 seconds"), h.Empty()),
		))
	}

	resource.RunStatus = domain.RunStatusRunning

	return h.SwapPartial(ctx, ResourceStatusContainer(resource))
}

func StopResource(ctx *h.RequestContext) *h.Partial {
	id := ctx.QueryParam("id")
	resource, err := resources.Stop(ctx.ServiceLocator(), id)

	if err != nil {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf(err.Error()), h.Empty()),
		))
	}

	// wait until its running
	m := service.Get[monitor.Monitor](ctx.ServiceLocator())

	success := util.WaitFor(time.Second*5, func() bool {
		return m.GetRunStatus(resource) != domain.RunStatusRunning
	})

	if !success {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf("Failed to stop resource after 5 seconds"), h.Empty()),
		))
	}

	resource.RunStatus = domain.RunStatusNotRunning

	return h.SwapPartial(ctx, ResourceStatusContainer(resource))
}
