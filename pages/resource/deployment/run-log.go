package deployment

import (
	"context"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
	"paas/app"
	"paas/pages/resource/resourceui"
	"paas/ui"
)

func RunLog(ctx *h.RequestContext) *h.Page {
	id := ctx.QueryParam("id")
	resource, err := app.ResourceGet(ctx.ServiceLocator(), id)

	if err != nil {
		ctx.Redirect("/", 302)
		return h.EmptyPage()
	}

	app.OnceWithAliveContext(ctx, func(context context.Context) {

		app.StreamLogs(ctx.ServiceLocator(), context, resource, func(msg *nats.Msg) {
			data := string(msg.Data)
			ws.PushElementCtx(ctx, ui.LogLine(data))
		})
	})

	return resourceui.Page(ctx, func(resource *app.Resource) *h.Element {
		return h.Div(
			h.Class("h-[500px]"),
			ui.LogBody(ui.LogBodyOptions{
				MaxLogs: 1000,
			}),
		)
	})
}
