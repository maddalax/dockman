package kv

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"paas/kv/subject"
)

func (c *Client) LogBuildError(resourceId string, buildId string, error error) {
	_, _ = c.js.Publish(subject.BuildLogForResource(resourceId, buildId), []byte(error.Error()))
}

func (c *Client) LogBuildMessage(resourceId string, buildId string, message string) {
	_, _ = c.js.Publish(subject.BuildLogForResource(resourceId, buildId), []byte(message))
}

func (c *Client) BuildLogStreamName(resourceId string, buildId string) string {
	return fmt.Sprintf("BUILD_LOG_STREAM-%s-%s", resourceId, buildId)
}

func (c *Client) CreateBuildLogStream(resourceId string, buildId string) error {
	_, err := c.js.AddStream(&nats.StreamConfig{
		Name: c.BuildLogStreamName(resourceId, buildId),
		// TODO should this have max age, and max msgs?
		Subjects:  []string{subject.BuildLogForResource(resourceId, buildId)},
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

func (c *Client) CreateHistoryStream() error {
	_, err := c.js.AddStream(&nats.StreamConfig{
		Name:      "HISTORY_STREAM",
		Subjects:  []string{string(subject.ResourceCreated)},
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
