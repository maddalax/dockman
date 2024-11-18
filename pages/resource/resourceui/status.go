package resourceui

import (
	"errors"
	"github.com/google/uuid"
	"github.com/maddalax/htmgo/framework/h"
	"paas/domain"
	"paas/resources"
	"paas/ui"
	"paas/urls"
)

func GetStatusPartial(ctx *h.RequestContext) *h.Partial {
	return resources.WithStatusLock(ctx.ServiceLocator(), ctx.QueryParam("id"), func(err error) *h.Partial {
		id := ctx.QueryParam("id")
		resource, err := resources.Get(ctx.ServiceLocator(), id)

		if err != nil {
			// TODO
			panic(err)
		}

		return h.NewPartial(
			PageHeader(ctx, resource),
		)
	})
}

func StartResource(ctx *h.RequestContext) *h.Partial {
	id := ctx.QueryParam("id")
	resource, err := resources.Start(ctx.ServiceLocator(), id, resources.StartOpts{
		IgnoreIfRunning: false,
	})

	if err != nil {
		// resource just hasn't been built yet, lets build it instead
		if errors.Is(err, domain.ResourceNotFoundError) {
			return h.RedirectPartial(urls.ResourceStartDeploymentPath(id, uuid.NewString()))
		}

		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf(err.Error()), h.Empty()),
		))
	}

	return h.SwapPartial(ctx, PageHeader(ctx, resource))
}

func RestartResource(ctx *h.RequestContext) *h.Partial {
	// start already handles the case where the resource is already running
	return StartResource(ctx)
}

func StopResource(ctx *h.RequestContext) *h.Partial {
	id := ctx.QueryParam("id")
	resource, err := resources.Stop(ctx.ServiceLocator(), id)

	if err != nil {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf(err.Error()), h.Empty()),
		))
	}

	return h.SwapPartial(ctx, PageHeader(ctx, resource))
}
