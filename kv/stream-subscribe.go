package kv

import (
	"context"
	"github.com/nats-io/nats.go"
	"io"
	"paas/kv/subject"
)

type WriterSubscriber struct {
	Writer     io.Writer
	Subscriber chan *nats.Msg
}

func (c *Client) CreateEphemeralWriterSubscriber(ctx context.Context, subject subject.Subject) *WriterSubscriber {
	ch := make(chan *nats.Msg, 100)
	_, err := c.SubscribeSubject(ctx, subject, func(msg *nats.Msg) {
		ch <- msg
	})
	if err != nil {
		return &WriterSubscriber{
			Writer:     &EmptyWriter{},
			Subscriber: ch,
		}
	}
	return &WriterSubscriber{
		Writer:     c.NewEphemeralNatsWriter(subject),
		Subscriber: ch,
	}
}
