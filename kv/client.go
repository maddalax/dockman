package kv

import (
	"errors"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
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
	client := service.Get[Client](ctx.ServiceLocator())
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
func (c *Client) SubscribeStreamAndReplayAll(subject subject.Subject, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := c.js.Subscribe(subject, handler, nats.DeliverAll())

	if err != nil {
		return nil, err
	}
	return sub, nil
}

// SubscribeStreamUntilTimeout subscribes to a subject and replays all messages, closing the subscription after no messages are received for the specified timeout
func (c *Client) SubscribeStreamUntilTimeout(subject subject.Subject, timeout time.Duration, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
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
			if time.Since(lastMessageTime) > timeout && sub != nil {
				_ = sub.Unsubscribe()
				return
			}
		}
	}()

	sub, err := c.js.Subscribe(subject, handle, nats.DeliverAll())
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (c *Client) SubscribeSubject(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := c.nc.Subscribe(subject, handler)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (c *Client) SubscribeStream(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := c.js.Subscribe(subject, handler, nats.StartTime(time.Now()))
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (c *Client) GetBucketWithConfig(config *nats.KeyValueConfig) (nats.KeyValue, error) {
	b, err := c.js.KeyValue(config.Bucket)
	if err != nil {
		if errors.Is(err, nats.ErrBucketNotFound) {
			b, err = c.js.CreateKeyValue(config)
			if err != nil {
				return nil, err
			}
			return b, nil
		}
	}
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
