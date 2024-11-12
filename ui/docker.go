package ui

import (
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
	"paas/builder"
	"paas/kv"
	"paas/kv/subject"
	"paas/resources"
	"paas/urls"
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
		h.Div(
			h.Class("flex gap-2"),
			PrimaryButton(ButtonProps{
				Text: "Re-run Build",
				Href: urls.ResourceStartDeploymentPath(resource.Id, buildId),
			}),
			PrimaryButton(ButtonProps{
				Text: "Cancel Build",
				Children: []h.Ren{
					ws.OnClick(ctx, func(data ws.HandlerData) {
						b := builder.GetBuilder(resource.Id, buildId)
						if b != nil {
							b.CancelBuild()
						}
					}),
				},
			}),
		),
		h.Div(
			h.Class("h-[500px]"),
			LogBody(),
		),
	)
}
