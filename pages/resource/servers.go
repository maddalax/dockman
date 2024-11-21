package resource

import (
	"dockside/app"
	"dockside/pages/resource/resourceui"
	"dockside/pages/servers"
	"github.com/maddalax/htmgo/framework/h"
)

func ServerPage(ctx *h.RequestContext) *h.Page {
	return resourceui.Page(ctx, func(resource *app.Resource) *h.Element {
		locator := ctx.ServiceLocator()

		details := h.Map(resource.ServerDetails, func(details app.ResourceServer) app.ResourceServerWithDetails {
			server, err := app.ServerGet(locator, details.ServerId)
			// this can happen if the server no longer exists and it hasn't been deleted from the resource yet
			if err != nil {
				return app.ResourceServerWithDetails{}
			}
			return app.ResourceServerWithDetails{
				ResourceServer: &details,
				Details:        server,
			}
		})

		return h.Div(
			h.Class("flex flex-col gap-4"),
			h.Name("resource-servers"),
			h.List(
				details,
				func(server app.ResourceServerWithDetails, index int) *h.Element {
					if server.Details != nil {
						return resourceServerBlock(&server)
					}
					return h.Empty()
				},
			),
		)
	})
}

func resourceServerBlock(server *app.ResourceServerWithDetails) *h.Element {
	return h.Div(
		servers.ServerBlockDetails(server.Details),
	)
}
