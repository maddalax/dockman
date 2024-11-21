package app

import (
	"dockside/app/subject"
	"github.com/nats-io/nats.go"
)

// NatsWriter is a structure that implements io.Writer to write to a NATS JetStream stream
type NatsWriter struct {
	js      nats.JetStreamContext
	closed  bool
	subject subject.Subject
}

func (c *KvClient) NewNatsWriter(subject subject.Subject) *NatsWriter {
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

func (nw *NatsWriter) Close() error {
	nw.closed = true
	return nil
}

type EphemeralNatsWriter struct {
	subject     subject.Subject
	closed      bool
	beforeWrite func(data string) bool
	c           *KvClient
}

func (c *KvClient) NewEphemeralNatsWriter(subject subject.Subject) *EphemeralNatsWriter {
	return &EphemeralNatsWriter{
		subject: subject,
		c:       c,
	}
}

func (nw *EphemeralNatsWriter) SetBeforeWrite(beforeWrite func(data string) bool) {
	nw.beforeWrite = beforeWrite
}

func (nw *EphemeralNatsWriter) IsClosed() bool {
	return nw.closed
}

func (nw *EphemeralNatsWriter) Close() error {
	nw.closed = true
	return nil
}

// Write implements the io.Writer interface
func (nw *EphemeralNatsWriter) Write(p []byte) (n int, err error) {
	if nw.beforeWrite != nil {
		shouldWrite := nw.beforeWrite(string(p))
		if !shouldWrite {
			return len(p), nil
		}
	}

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

func (ew *EmptyWriter) Close() error {
	return nil
}
