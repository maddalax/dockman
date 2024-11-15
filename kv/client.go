package kv

import (
	"context"
	"errors"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"log/slog"
	"paas/kv/subject"
	"time"
)

type Options struct {
	Port int
}

type Client struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func GetClientFromCtx(ctx *h.RequestContext) *Client {
	return GetClientFromLocator(ctx.ServiceLocator())
}

func GetClientFromLocator(locator *service.Locator) *Client {
	client := service.Get[Client](locator)
	return client
}

func (c *Client) PurgeStream(stream string) error {
	return c.js.PurgeStream(stream)
}

func (c *Client) DeleteStream(stream string) error {
	return c.js.DeleteStream(stream)
}

func (c *Client) Publish(subject string, data []byte) error {
	return c.nc.Publish(subject, data)
}

func (c *Client) DeleteBucket(bucket string) error {
	return c.js.DeleteKeyValue(bucket)
}

// SubscribeStreamAndReplayAll subscribes to a stream and replays all messages
func (c *Client) SubscribeStreamAndReplayAll(context context.Context, subject subject.Subject, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := c.js.Subscribe(subject, handler, nats.DeliverAll())

	if err != nil {
		return nil, err
	}

	DisposeOnCancel(context, func() {
		_ = sub.Unsubscribe()
	})

	return sub, nil
}

// SubscribeStreamUntilTimeout subscribes to a subject and replays all messages, closing the subscription after no messages are received for the specified timeout
func (c *Client) SubscribeStreamUntilTimeout(context context.Context, subject subject.Subject, timeout time.Duration, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	lastMessageTime := time.Now()
	var sub *nats.Subscription

	handle := func(msg *nats.Msg) {
		handler(msg)
		lastMessageTime = time.Now()
		sub = msg.Sub
	}

	go func() {
		for {
			time.Sleep(time.Second)
			select {
			case <-context.Done():
				if sub != nil {
					slog.Debug("Unsubscribing from stream", slog.String("subject", subject))
					_ = sub.Unsubscribe()
					return
				}
			default:
				if time.Since(lastMessageTime) > timeout && sub != nil {
					slog.Debug("Unsubscribing from stream", slog.String("subject", subject))
					_ = sub.Unsubscribe()
					return
				}
			}
		}
	}()

	slog.Debug("Subscribing to stream", slog.String("subject", subject))
	sub, err := c.js.Subscribe(subject, handle, nats.DeliverAll())
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (c *Client) SubscribeSubject(context context.Context, subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := c.nc.Subscribe(subject, handler)
	if err != nil {
		return nil, err
	}
	DisposeOnCancel(context, func() {
		_ = sub.Unsubscribe()
	})
	return sub, nil
}

func (c *Client) SubscribeStream(context context.Context, subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := c.js.Subscribe(subject, handler, nats.StartTime(time.Now()))
	if err != nil {
		return nil, err
	}
	DisposeOnCancel(context, func() {
		_ = sub.Unsubscribe()
	})
	return sub, nil
}

func (c *Client) GetBucketWithConfig(config *nats.KeyValueConfig) (nats.KeyValue, error) {
	b, err := c.js.KeyValue(config.Bucket)
	return b, err
}

func (c *Client) GetStreams() []*nats.StreamInfo {
	stores := c.js.Streams()
	streams := make([]*nats.StreamInfo, 0)
	for store := range stores {
		streams = append(streams, store)
	}
	return streams
}

func (c *Client) GetBuckets() []nats.KeyValueStatus {
	stores := c.js.KeyValueStores()
	buckets := make([]nats.KeyValueStatus, 0)
	for store := range stores {
		buckets = append(buckets, store)
	}
	return buckets
}

func (c *Client) GetStream(bucket string) (*nats.StreamInfo, error) {
	return c.js.StreamInfo(bucket)
}

func (c *Client) GetBucket(bucket string) (nats.KeyValue, error) {
	return c.GetBucketWithConfig(&nats.KeyValueConfig{
		Bucket: bucket,
	})
}

func (c *Client) GetOrCreateBucket(bucket string) (nats.KeyValue, error) {
	b, err := c.GetBucketWithConfig(&nats.KeyValueConfig{
		Bucket: bucket,
	})
	if err != nil {
		if errors.Is(err, nats.ErrBucketNotFound) {
			b, err = c.CreateBucket(&nats.KeyValueConfig{
				Bucket: bucket,
			})
			return b, err
		}
	}
	return b, err
}

func (c *Client) CreateBucket(config *nats.KeyValueConfig) (nats.KeyValue, error) {
	return c.js.CreateKeyValue(config)
}

func Connect(opts Options) (*Client, error) {
	// Connect to the embedded NATS server
	nc, err := nats.Connect(fmt.Sprintf("nats://localhost:%d", opts.Port))
	if err != nil {
		return nil, err
	}

	// Use JetStream
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &Client{
		nc: nc,
		js: js,
	}, nil
}
