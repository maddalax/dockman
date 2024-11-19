package jetstream

import (
	"context"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
	"paas/internal"
	"paas/internal/ui"
	"paas/pages"
	"time"
)

func StreamsDebugPage(ctx *h.RequestContext) *h.Page {
	return pages.RootPage(
		ctx,
		h.Div(
			h.Class("flex flex-row min-h-screen"),
			StreamsSidebar(ctx),
			h.Div(
				h.Class("flex flex-col gap-4 items-center w-3/4 pt-8"),
				h.Div(
					h.Class("mt-3"),
					h.H3(
						h.Text("JetStream Debug Page"),
						h.Class("text-xl font-bold text-center mb-4"),
					),
					h.Div(
						h.Id("key-value-pairs"),
						h.Class("flex flex-col gap-2 items-center mt-4 w-full"),
					),
				),
			),
		),
	)
}

func StreamsSidebar(ctx *h.RequestContext) *h.Element {
	client := internal.KvFromCtx(ctx)
	buckets := client.GetStreams()
	return h.Div(
		h.Id("bucket-list"),
		h.Class("w-1/4 flex flex-col gap-2 items-start p-2 bg-gray-200 overflow-y-auto bg-neutral-50 px-4"),
		h.H4(
			h.Text("Buckets"),
			h.Class("font-bold mb-2"),
		),
		h.List(
			buckets,
			func(bucket *nats.StreamInfo, index int) *h.Element {
				return StreamCard(ctx, bucket)
			},
		),
	)
}

func StreamCard(ctx *h.RequestContext, streamStatus *nats.StreamInfo) *h.Element {
	client := internal.KvFromCtx(ctx)
	deleteButton := h.Button(
		h.Class("text-blue underline"),
		h.Text("Delete"),
		ws.OnClick(ctx, func(data ws.HandlerData) {
			client.DeleteStream(streamStatus.Config.Name)
			ws.PushElement(data, EmptyDetails())
			ws.PushElement(data, StreamsSidebar(ctx))
		}),
	)

	return h.Div(
		h.Class("flex flex-row gap-3 items-center w-full"),
		h.Div(
			h.Class("flex flex-col gap-1 border-r border-slate-200 w-full"),
			ws.OnClick(ctx, func(data ws.HandlerData) {
				ws.PushElementCtx(ctx, StreamDetails(ctx, streamStatus))
			}),
			h.Pf(
				streamStatus.Config.Name,
			),
		),
		deleteButton,
	)
}

func StreamDetails(ctx *h.RequestContext, stream *nats.StreamInfo) *h.Element {

	client := internal.KvFromCtx(ctx)

	internal.OnceWithAliveContext(ctx, func(context context.Context) {
		for _, subject := range stream.Config.Subjects {
			client.SubscribeStreamUntilTimeout(context, subject, time.Second*3, func(msg *nats.Msg) {
				data := string(msg.Data)
				ws.PushElementCtx(ctx, ui.LogLine(data))
			})
		}
	})

	return h.Div(
		h.Id("key-value-pairs"),
		h.Class("flex flex-col gap-4 items-center mt-4 w-full p-4 border border-slate-200 rounded-md"),
		h.H4(
			h.Text(stream.Config.Name),
			h.Class("font-bold"),
		),

		h.Div(
			h.Class("flex h-[800px]"),
			ui.LogBody(ui.LogBodyOptions{
				MaxLogs: 100,
			}),
		),
	)
}
