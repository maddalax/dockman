package kv

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type Options struct {
	Port int
}

type Client struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func (c *Client) Publish(subject string, data []byte) error {
	return c.nc.Publish(subject, data)
}

// SubscribeAndReplayAll subscribes to a subject and replays all messages
func (c *Client) SubscribeAndReplayAll(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := c.js.Subscribe(subject, handler, nats.DeliverAll())
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (c *Client) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
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
	return b, nil
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
