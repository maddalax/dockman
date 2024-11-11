package jetstream

import (
	"fmt"
	"github.com/maddalax/htmgo/extensions/websocket/session"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
	"paas/kv"
	"paas/pages"
)

func KvDebugPage(ctx *h.RequestContext) *h.Page {
	sessionId := session.GetSessionId(ctx)
	return pages.RootPage(
		h.Div(
			h.Attribute("ws-connect", fmt.Sprintf("/ws?sessionId=%s", sessionId)),
			h.Class("flex flex-row min-h-screen"),
			BucketSidebar(ctx),
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

func EmptyDetails() *h.Element {
	return h.Div(
		h.Id("key-value-pairs"),
		h.Class("flex flex-col gap-2 items-center mt-4 w-full"),
	)
}

func BucketSidebar(ctx *h.RequestContext) *h.Element {
	client := kv.GetClientFromCtx(ctx)
	buckets := client.GetBuckets()
	return h.Div(
		h.Id("bucket-list"),
		h.Class("w-1/4 flex flex-col gap-2 items-start p-2 bg-gray-200 overflow-y-auto bg-neutral-50 px-4"),
		h.H4(
			h.Text("Buckets"),
			h.Class("font-bold mb-2"),
		),
		h.List(
			buckets,
			func(bucket nats.KeyValueStatus, index int) *h.Element {
				return BucketCard(ctx, bucket)
			},
		),
	)
}

func BucketCard(ctx *h.RequestContext, bucketStatus nats.KeyValueStatus) *h.Element {
	client := kv.GetClientFromCtx(ctx)
	deleteButton := h.Button(
		h.Class("text-blue underline"),
		h.Text("Delete"),
		ws.OnClick(ctx, func(data ws.HandlerData) {
			client.DeleteBucket(bucketStatus.Bucket())
			ws.PushElement(data, EmptyDetails())
			ws.PushElement(data, BucketSidebar(ctx))
		}),
	)

	return h.Div(
		h.Class("flex flex-row gap-3 items-center w-full"),
		h.Div(
			h.Class("flex flex-col gap-1 border-r border-slate-200 w-full"),
			ws.OnClick(ctx, func(data ws.HandlerData) {
				ws.PushElementCtx(ctx, BucketDetails(ctx, bucketStatus))
			}),
			h.Pf(
				bucketStatus.Bucket(),
			),
		),
		deleteButton,
	)
}

func BucketDetails(ctx *h.RequestContext, bucketStatus nats.KeyValueStatus) *h.Element {
	client := kv.GetClientFromCtx(ctx)
	bucket, err := client.GetBucket(bucketStatus.Bucket())

	if err != nil {
		return h.Div()
	}

	keyChan, _ := bucket.ListKeys()
	keys := make([]string, 0)

	for key := range keyChan.Keys() {
		keys = append(keys, key)
	}

	return h.Div(
		h.Id("key-value-pairs"),
		h.Class("flex flex-col gap-4 items-center mt-4 w-full p-4 border border-slate-200 rounded-md"),
		h.H4(
			h.Text(bucketStatus.Bucket()),
			h.Class("font-bold"),
		),
		h.List(
			keys,
			func(key string, index int) *h.Element {
				value, _ := bucket.Get(key)

				deleteButton := h.Button(
					h.Class("text-blue underline"),
					h.Text("Delete"),
					ws.OnClick(ctx, func(data ws.HandlerData) {
						bucket.Delete(key)
						ws.PushElement(data, BucketDetails(ctx, bucketStatus))
					}),
				)

				return h.Div(
					h.Class("flex flex-row gap-3 items-center w-full"),
					h.Div(
						h.Class("flex flex-col gap-1 p-2 border border-slate-200 rounded-md w-full"),
						h.Span(
							h.Class("font-bold"),
							h.Text(key),
						),
						h.Span(
							h.Text(string(value.Value())),
						),
					),
					deleteButton,
				)
			},
		),
	)
}
