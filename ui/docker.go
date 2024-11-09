package ui

import (
	"github.com/docker/docker/api/types"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
	"github.com/nats-io/nats.go"
	"paas/docker"
	"paas/kv"
)

func DockerBuildTest(ctx *h.RequestContext) *h.Element {
	natsClient, err := kv.Connect(kv.Options{
		Port: 4222,
	})
	if err != nil {
		panic(err)
	}

	natsClient.CreateBuildLogStream()

	ws.Once(ctx, func() {
		natsClient.SubscribeAndReplayAll(kv.BuildLogStreamSubject, func(msg *nats.Msg) {
			data := string(msg.Data)
			ws.PushElementCtx(ctx, h.Div(
				h.Class("bg-slate-50 p-2 rounded-md text-sm"),
				h.Attribute("hx-swap-oob", "beforeend:#build-log"),
				h.P(
					h.Text(data),
				),
			))
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
				outputStream := natsClient.NewNatsWriter(kv.BuildLogStreamSubject)
				client.Build(outputStream, ".", types.ImageBuildOptions{
					Dockerfile: "Dockerfile",
				})
			}),
		),
		h.Div(
			h.Class("max-w-[800px] max-h-80 overflow-y-auto bg-white border border-gray-300 rounded-lg shadow-lg p-4 mt-6 min-w-[800px]"),
			h.Id("build-log"),
			// Scroll to the bottom of the div when the page loads
			h.OnLoad(
				// language=JavaScript
				js.EvalJs(`
					setTimeout(() => {
           self.scrollTop = self.scrollHeight;       
					}, 1000)
				`),
			),
			// Scroll to the bottom of the div when the message is sent
			// only if the user is close to the bottom of the div
			h.OnEvent("htmx:wsAfterMessage",
				// language=JavaScript
				js.EvalJs(`
					const scrollPosition = self.scrollTop + self.clientHeight;
    			const distanceFromBottom = self.scrollHeight - scrollPosition;
    			const scrollThreshold = 1000; // Adjust this to define how close the user should be to the bottom
					 if (distanceFromBottom <= scrollThreshold) {
        			self.scrollTop = self.scrollHeight;
    			}
				`),
			),
		),
	)
}
