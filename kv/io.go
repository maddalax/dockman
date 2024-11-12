package kv

import (
	"github.com/nats-io/nats.go"
	"paas/kv/subject"
)

// NatsWriter is a structure that implements io.Writer to write to a NATS JetStream stream
type NatsWriter struct {
	js      nats.JetStreamContext
	subject subject.Subject
}

func (c *Client) NewNatsWriter(subject subject.Subject) *NatsWriter {
	return &NatsWriter{
		js:      c.js,
		subject: subject,
	}
}

// Write implements the io.Writer interface
func (nw *NatsWriter) Write(p []byte) (n int, err error) {
	_, err = nw.js.Publish(nw.subject, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

type EphemeralNatsWriter struct {
	subject subject.Subject
	c       *Client
}

func (c *Client) NewEphemeralNatsWriter(subject subject.Subject) *EphemeralNatsWriter {
	return &EphemeralNatsWriter{
		subject: subject,
		c:       c,
	}
}

// Write implements the io.Writer interface
func (nw *EphemeralNatsWriter) Write(p []byte) (n int, err error) {
	err = nw.c.Publish(nw.subject, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

type EmptyWriter struct {
}

func (ew *EmptyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
