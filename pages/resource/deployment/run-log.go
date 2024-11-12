package deployment

import (
	"fmt"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
	"paas/docker"
	"paas/kv"
	"paas/kv/subject"
	"paas/pages"
	"paas/resources"
	"paas/ui"
)

func RunLog(ctx *h.RequestContext) *h.Page {
	id := ctx.QueryParam("id")
	resource, err := resources.Get(ctx.ServiceLocator(), id)

	if err != nil {
		ctx.Redirect("/", 302)
		return h.EmptyPage()
	}

	natsClient := kv.GetClientFromCtx(ctx)

	ws.Once(ctx, func() {
		natsClient.SubscribeSubject(subject.RunLogsForResource(resource.Id), func(msg *nats.Msg) {
			data := string(msg.Data)
			data = "RUN: " + data
			ws.PushElementCtx(ctx, ui.LogLine(data))
		})
		writer := natsClient.NewEphemeralNatsWriter(
			subject.RunLogsForResource(resource.Id),
		)
		client, err := docker.Connect()
		if err != nil {
			ws.PushElementCtx(ctx, ui.LogLine("Failed connecting to docker"))
		} else {
			containerId := fmt.Sprintf("%s-%s-container", resource.Name, resource.Id)
			if err == nil {
				err := client.StreamLogs(containerId, docker.StreamLogsOptions{
					Stdout: writer,
				})
				if err != nil {
					ws.PushElementCtx(ctx, ui.LogLine("Failed to stream logs"))
				}
			}
		}
	})

	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-2"),
			pages.Title("Run Logs"),
			h.Pf("Resource: %s", resource.Name),
			TopTabs(ctx, resource),
			h.Div(
				h.Class("h-[500px]"),
				ui.LogBody(),
			),
		),
	)
}
