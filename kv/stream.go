package kv

import (
	"github.com/nats-io/nats.go"
)

var BuildLogStreamSubject = "build.log"

func (c *Client) CreateBuildLogStream() error {
	_, err := c.js.AddStream(&nats.StreamConfig{
		Name:      "BUILD_LOG_STREAM",
		Subjects:  []string{BuildLogStreamSubject},
		Retention: nats.LimitsPolicy, // Retain messages until storage limit is reached
		MaxAge:    0,                 // Messages never expire based on age
		MaxMsgs:   -1,                // No limit on the number of messages
		MaxBytes:  -1,                // No limit on the total size of messages
		Storage:   nats.FileStorage,  // Use file storage for persistence
	})
	if err != nil {
		return err
	}
	return nil
}
