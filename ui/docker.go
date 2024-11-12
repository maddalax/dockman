package ui

import (
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
	"paas/kv"
	"paas/kv/subject"
	"paas/resources"
)

func DockerBuildLogs(ctx *h.RequestContext, resource *resources.Resource, buildId string) *h.Element {
	natsClient := kv.GetClientFromCtx(ctx)
	ws.Once(ctx, func() {
		natsClient.SubscribeAndReplayAll(subject.BuildLogForResource(resource.Id, buildId), func(msg *nats.Msg) {
			data := string(msg.Data)
			ws.PushElementCtx(ctx, LogLine(data))
		})
	})

	return h.Div(
		h.Class("flex flex-col gap-6 items-center justify-center p-8"),
		h.H2(
			h.Text("Deployment Log"),
			h.Class("text-xl font-bold text-center"),
		),
		LogBody(),
	)
}
