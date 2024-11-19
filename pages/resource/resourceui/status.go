package resourceui

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/app"
	"paas/app/ui"
	"time"
)

func GetStatusPartial(ctx *h.RequestContext) *h.Partial {
	return app.WithStatusLock(ctx.ServiceLocator(), ctx.QueryParam("id"), func(err error) *h.Partial {
		id := ctx.QueryParam("id")
		resource, err := app.ResourceGet(ctx.ServiceLocator(), id)

		if err != nil {
			// TODO
			panic(err)
		}

		return h.SwapPartial(
			ctx,
			PageHeader(ctx, resource),
		)
	})
}

func StartResource(ctx *h.RequestContext) *h.Partial {
	id := ctx.QueryParam("id")

	_, err := app.SendCommand[app.RunResourceResponse](ctx.ServiceLocator(), app.SendCommandOpts{
		Command: &app.RunResourceCommand{
			ResourceId: id,
		},
		Timeout: time.Second * 5,
	})

	if err != nil {
		//// resource just hasn't been built yet, lets build it instead
		//if errors.Is(err, internal.ResourceNotFoundError) {
		//	return h.RedirectPartial(urls.ResourceStartDeploymentPath(id, uuid.NewString()))
		//}

		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf(err.Error()), h.Empty()),
		))
	}

	resource, err := app.ResourceGet(ctx.ServiceLocator(), id)

	return h.SwapPartial(ctx, PageHeader(ctx, resource))
}

func RestartResource(ctx *h.RequestContext) *h.Partial {
	// start already handles the case where the resource is already running
	return StartResource(ctx)
}

func StopResource(ctx *h.RequestContext) *h.Partial {
	id := ctx.QueryParam("id")

	_, err := app.SendCommand[app.StopResourceResponse](ctx.ServiceLocator(), app.SendCommandOpts{
		Command: &app.StopResourceCommand{
			ResourceId: id,
		},
		Timeout: time.Second * 5,
	})

	if err != nil {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf(err.Error()), h.Empty()),
		))
	}

	resource, err := app.ResourceGet(ctx.ServiceLocator(), id)

	return h.SwapPartial(ctx, PageHeader(ctx, resource))
}
