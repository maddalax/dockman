package app

import (
	"context"
	"dockside/app/logger"
	"dockside/app/subject"
	"errors"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"log"
	"net"
	"time"
)

type NatsConnectOptions struct {
	Port int
	Host string
}

type KvClient struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func KvFromCtx(ctx *h.RequestContext) *KvClient {
	return KvFromLocator(ctx.ServiceLocator())
}

func KvFromLocator(locator *service.Locator) *KvClient {
	client := service.Get[KvClient](locator)
	return client
}

func (c *KvClient) Ping() error {
	return c.nc.Flush()
}

func (c *KvClient) PurgeStream(stream string) error {
	return c.js.PurgeStream(stream)
}

func (c *KvClient) DeleteStream(stream string) error {
	return c.js.DeleteStream(stream)
}

func (c *KvClient) Publish(subject string, data []byte) error {
	return c.nc.Publish(subject, data)
}

func (c *KvClient) DeleteBucket(bucket string) error {
	return c.js.DeleteKeyValue(bucket)
}

// SubscribeStreamAndReplayAll subscribes to a stream and replays all messages
func (c *KvClient) SubscribeStreamAndReplayAll(context context.Context, subject subject.Subject, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
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
func (c *KvClient) SubscribeStreamUntilTimeout(context context.Context, subject subject.Subject, timeout time.Duration, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
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
					logger.DebugWithFields("Unsubscribing from stream", map[string]any{
						"subject": subject,
					})
					_ = sub.Unsubscribe()
					return
				}
			default:
				if time.Since(lastMessageTime) > timeout && sub != nil {
					logger.DebugWithFields("Unsubscribing from stream due to timeout", map[string]any{
						"subject": subject,
					})
					_ = sub.Unsubscribe()
					return
				}
			}
		}
	}()

	logger.DebugWithFields("Subscribing to stream", map[string]any{
		"subject": subject,
	})
	sub, err := c.js.Subscribe(subject, handle, nats.DeliverAll())
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (c *KvClient) SubscribeSubject(context context.Context, subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := c.nc.Subscribe(subject, handler)
	if err != nil {
		return nil, err
	}
	DisposeOnCancel(context, func() {
		_ = sub.Unsubscribe()
	})
	return sub, nil
}

func (c *KvClient) SubscribeSubjectForever(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg)
	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (c *KvClient) SubscribeStream(context context.Context, subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := c.js.Subscribe(subject, handler, nats.StartTime(time.Now()))
	if err != nil {
		return nil, err
	}
	DisposeOnCancel(context, func() {
		_ = sub.Unsubscribe()
	})
	return sub, nil
}

func (c *KvClient) GetBucketWithConfig(config *nats.KeyValueConfig) (nats.KeyValue, error) {
	b, err := c.js.KeyValue(config.Bucket)
	return b, err
}

func (c *KvClient) GetStreams() []*nats.StreamInfo {
	stores := c.js.Streams()
	streams := make([]*nats.StreamInfo, 0)
	for store := range stores {
		streams = append(streams, store)
	}
	return streams
}

func (c *KvClient) GetBuckets() []nats.KeyValueStatus {
	stores := c.js.KeyValueStores()
	buckets := make([]nats.KeyValueStatus, 0)
	for store := range stores {
		buckets = append(buckets, store)
	}
	return buckets
}

func (c *KvClient) GetStream(bucket string) (*nats.StreamInfo, error) {
	return c.js.StreamInfo(bucket)
}

func (c *KvClient) GetBucket(bucket string) (nats.KeyValue, error) {
	return c.GetBucketWithConfig(&nats.KeyValueConfig{
		Bucket: bucket,
	})
}

func (c *KvClient) GetOrCreateBucket(config *nats.KeyValueConfig) (nats.KeyValue, error) {
	b, err := c.GetBucketWithConfig(config)
	if err != nil {
		if errors.Is(err, nats.ErrBucketNotFound) {
			b, err = c.CreateBucket(config)
			return b, err
		}
	}
	return b, err
}

func (c *KvClient) CreateBucket(config *nats.KeyValueConfig) (nats.KeyValue, error) {
	return c.js.CreateKeyValue(config)
}

func NatsConnect(opts NatsConnectOptions) (*KvClient, error) {
	natsOpts := []nats.Option{
		nats.Name("Retry Connection Example"),
		nats.ReconnectWait(2 * time.Second),
		// always try to reconnect
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("Disconnected due to: %v, will attempt to reconnect...", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Reconnected to %v", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("Connection closed.")
		}),
	}

	if opts.Host == "" {
		opts.Host = "localhost"
	}

	host := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)

	nc, err := nats.Connect(host, natsOpts...)

	if err != nil {

		canReconnect := false

		var opError *net.OpError
		switch {
		case errors.As(err, &opError):
			canReconnect = true
		}
		if err.Error() == "nats: no servers available for connection" {
			canReconnect = true
		}

		if canReconnect {
			for {
				log.Printf("Retrying nats connection to %s", host)
				nc, err = nats.Connect(host, natsOpts...)
				if err == nil {
					break
				}
				time.Sleep(2 * time.Second)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	if nc == nil {
		return nil, errors.New("nats connection is nil")
	}

	// Use JetStream
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &KvClient{
		nc: nc,
		js: js,
	}, nil
}
