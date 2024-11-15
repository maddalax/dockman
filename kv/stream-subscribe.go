package kv

import (
	"context"
	"github.com/nats-io/nats.go"
	"io"
	"paas/kv/subject"
)

type WriterSubscriber struct {
	Writer     io.WriteCloser
	Subscriber chan *nats.Msg
}

type CreateOptions struct {
	BeforeWrite func(data string) bool
}

func (c *Client) CreateEphemeralWriterSubscriber(ctx context.Context, subject subject.Subject, opts CreateOptions) *WriterSubscriber {
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

	w := c.NewEphemeralNatsWriter(subject)
	if opts.BeforeWrite != nil {
		w.SetBeforeWrite(opts.BeforeWrite)
	}

	return &WriterSubscriber{
		Writer:     w,
		Subscriber: ch,
	}
}
