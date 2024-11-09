package kv

import (
	"github.com/nats-io/nats.go"
)

// NatsWriter is a structure that implements io.Writer to write to a NATS JetStream stream
type NatsWriter struct {
	js      nats.JetStreamContext
	subject string
}

func (c *Client) NewNatsWriter(subject string) *NatsWriter {
	return &NatsWriter{
		js:      c.js,
		subject: subject,
	}
}

// Write implements the io.Writer interface
func (nw *NatsWriter) Write(p []byte) (n int, err error) {
	// Publish the data to the NATS JetStream subject
	_, err = nw.js.Publish(nw.subject, p)
	if err != nil {
		return 0, err
	}
	// Return the length of the written data
	return len(p), nil
}
