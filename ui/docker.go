package ui

import (
	"github.com/docker/docker/api/types"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
	"paas/docker"
	"paas/kv"
	"paas/kv/subject"
	"time"
)

func DockerBuildTest(ctx *h.RequestContext) *h.Element {
	// todo make singleton client
	natsClient, err := kv.Connect(kv.Options{
		Port: 4222,
	})
	if err != nil {
		panic(err)
	}

	ws.Every(ctx, time.Second, func() bool {
		now := time.Now()
		natsClient.Publish(string(subject.BuildLog), []byte(now.Format(time.Stamp)))
		return true
	})

	ws.Once(ctx, func() {
		// todo move this to app entry
		natsClient.CreateBuildLogStream()
		natsClient.SubscribeAndReplayAll(subject.BuildLog, func(msg *nats.Msg) {
			data := string(msg.Data)
			ws.PushElementCtx(ctx, LogLine(data))
		})
	})

	return h.Div(
		h.Class("flex flex-col gap-6 items-center justify-center bg-gray-100 p-8"),
		h.H2(
			h.Text("Docker Build Log"),
			h.Class("text-3xl font-bold text-center mb-4"),
		),
		h.Button(
			h.Text("Start build"),
			h.Class("bg-rose-400 hover:bg-rose-500 text-white font-bold py-2 px-4 rounded"),
			ws.OnClick(ctx, func(data ws.HandlerData) {
				client, err := docker.Connect()
				if err != nil {
					panic(err)
				}
				outputStream := natsClient.NewNatsWriter(subject.BuildLog)
				client.Build(outputStream, ".", types.ImageBuildOptions{
					Dockerfile: "Dockerfile",
				})
			}),
		),
		LogBody(),
	)
}
