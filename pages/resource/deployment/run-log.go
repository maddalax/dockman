package deployment

import (
	"context"
	"fmt"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
	"paas/logger"
	"paas/pages"
	resource2 "paas/pages/resource"
	"paas/resources"
	"paas/ui"
	"paas/wsutil"
)

func RunLog(ctx *h.RequestContext) *h.Page {
	id := ctx.QueryParam("id")
	resource, err := resources.Get(ctx.ServiceLocator(), id)

	if err != nil {
		ctx.Redirect("/", 302)
		return h.EmptyPage()
	}

	wsutil.OnceWithAliveContext(ctx, func(context context.Context) {
		logger.StreamLogs(ctx.ServiceLocator(), context, resource, func(msg *nats.Msg) {
			data := string(msg.Data)
			data = "RUN: " + data
			success := ws.PushElementCtx(ctx, ui.LogLine(data))
			fmt.Printf("Pushed log line: %s, success: %v\n", data, success)
		})
	})

	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-2"),
			pages.Title("Run Logs"),
			h.Pf("Resource: %s", resource.Name),
			resource2.TopTabs(ctx, resource),
			h.Div(
				h.Class("h-[500px]"),
				ui.LogBody(),
			),
		),
	)
}
