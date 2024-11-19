package deployment

import (
	"context"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
	"paas/internal"
	"paas/internal/ui"
	"paas/pages/resource/resourceui"
)

func RunLog(ctx *h.RequestContext) *h.Page {
	id := ctx.QueryParam("id")
	resource, err := internal.ResourceGet(ctx.ServiceLocator(), id)

	if err != nil {
		ctx.Redirect("/", 302)
		return h.EmptyPage()
	}

	internal.OnceWithAliveContext(ctx, func(context context.Context) {

		internal.StreamResourceLogs(ctx.ServiceLocator(), context, resource, func(msg *nats.Msg) {
			data := string(msg.Data)
			ws.PushElementCtx(ctx, ui.LogLine(data))
		})
	})

	return resourceui.Page(ctx, func(resource *internal.Resource) *h.Element {
		return h.Div(
			h.Class("h-[500px]"),
			ui.LogBody(ui.LogBodyOptions{
				MaxLogs: 1000,
			}),
		)
	})
}
