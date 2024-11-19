package app

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"paas/app/subject"
)

func (c *KvClient) LogBuildError(resourceId string, buildId string, error error) {
	str := error.Error()
	str = fmt.Sprintf("BUILD_ERROR: %s", str)
	_, _ = c.js.Publish(subject.BuildLogForResource(resourceId, buildId), []byte(str))
}

func (c *KvClient) LogBuildMessage(resourceId string, buildId string, message string) {
	_, _ = c.js.Publish(subject.BuildLogForResource(resourceId, buildId), []byte(message))
}

func (c *KvClient) LogRunMessage(resourceId string, message string) {
	_, _ = c.js.Publish(subject.RunLogsForResource(resourceId), []byte(message))
}

func (c *KvClient) BuildLogStreamName(resourceId string, buildId string) string {
	return fmt.Sprintf("BUILD_LOG_STREAM-%s-%s", resourceId, buildId)
}

func (c *KvClient) RunLogStreamName(resourceId string) string {
	return fmt.Sprintf("RUN_LOG_STREAM-%s", resourceId)
}
func (c *KvClient) CreateRunLogStream(resourceId string) error {
	_, err := c.js.AddStream(&nats.StreamConfig{
		Name: c.RunLogStreamName(resourceId),
		// TODO should this have max age, and max msgs?
		Subjects:  []string{subject.RunLogsForResource(resourceId)},
		Retention: nats.LimitsPolicy, // Retain messages until storage limit is reached
		MaxAge:    0,                 // Messages never expire based on age
		MaxMsgs:   10 * 1000,         // No limit on the number of messages
		MaxBytes:  -1,                // No limit on the total size of messages
		Storage:   nats.FileStorage,  // Use file storage for persistence
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *KvClient) CreateBuildLogStream(resourceId string, buildId string) error {
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

func (c *KvClient) CreateHistoryStream() error {
	config := &nats.StreamConfig{
		Name: "HISTORY_STREAM",
		Subjects: []string{
			subject.ResourceCreated,
			subject.ResourcePatched,
			subject.ResourceStarted,
			subject.ResourceStopped,
		},
		Retention: nats.LimitsPolicy, // Retain messages until storage limit is reached
		MaxAge:    0,                 // Messages never expire based on age
		MaxMsgs:   -1,                // No limit on the number of messages
		MaxBytes:  -1,                // No limit on the total size of messages
		Storage:   nats.FileStorage,  // Use file storage for persistence
	}

	_, err := c.js.AddStream(config)
	if err != nil {
		var APIError *nats.APIError
		switch {
		case errors.As(err, &APIError):
			if APIError.ErrorCode == nats.JSErrCodeStreamNameInUse {
				// stream already exists, lets just update it
				_, err = c.js.UpdateStream(config)
			}
		}
		return err
	}
	return nil
}
