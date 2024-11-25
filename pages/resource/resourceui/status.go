package resourceui

import (
	"dockside/app"
	"dockside/app/ui"
	"github.com/maddalax/htmgo/framework/h"
)

// GetStatusPartial TODO update this to consult all the servers
func GetStatusPartial(ctx *h.RequestContext) *h.Partial {
	return app.WithStatusLock(ctx.ServiceLocator(), ctx.QueryParam("id"), func(err error) *h.Partial {
		id := ctx.QueryParam("id")
		resource, err := app.ResourceGet(ctx.ServiceLocator(), id)

		if err != nil {
			return ui.GenericErrorAlertPartial(ctx, err)
		}

		return h.SwapPartial(
			ctx,
			PageHeader(ctx, resource),
		)
	})
}

func StartResource(ctx *h.RequestContext) *h.Partial {
	id := ctx.QueryParam("id")

	_, err := app.SendResourceStartCommand(ctx.ServiceLocator(), id, app.StartOpts{
		RemoveExisting: true,
	})

	if err != nil {
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

	_, err := app.SendResourceStopCommand(ctx.ServiceLocator(), id)

	if err != nil {
		return h.SwapPartial(ctx, h.Fragment(
			ui.ErrorAlert(h.Pf(err.Error()), h.Empty()),
		))
	}

	resource, err := app.ResourceGet(ctx.ServiceLocator(), id)

	return h.SwapPartial(ctx, PageHeader(ctx, resource))
}
