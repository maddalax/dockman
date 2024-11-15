package resourceui

import (
	"github.com/maddalax/htmgo/framework/h"
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
	resource, err := resources.Start(ctx.ServiceLocator(), id)

	if err != nil {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf(err.Error()), h.Empty()),
		))
	}

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

	return h.SwapPartial(ctx, ResourceStatusContainer(resource))
}
