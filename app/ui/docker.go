package ui

import (
	"context"
	"dockside/app"
	"dockside/app/subject"
	"dockside/app/urls"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
)

func DockerBuildLogs(ctx *h.RequestContext, resource *app.Resource, buildId string) *h.Element {
	natsClient := app.KvFromCtx(ctx)

	app.OnceWithAliveContext(ctx, func(context context.Context) {
		sb := subject.BuildLogForResource(resource.Id, buildId)
		natsClient.SubscribeStreamAndReplayAll(context, sb, func(msg *nats.Msg) {
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
						registry := app.GetBuilderRegistry(ctx.ServiceLocator())
						b := registry.GetBuilder(resource.Id, buildId)
						if b != nil {
							b.CancelBuild()
						}
					}),
				},
			}),
		),
		h.Div(
			h.Class("h-[calc(100vh-400px)] w-full"),
			LogBody(
				LogBodyOptions{
					MaxLogs: 1000,
				},
			),
		),
	)
}
