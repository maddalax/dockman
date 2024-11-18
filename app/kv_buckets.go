package app

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

func (c *KvClient) GetResourceDeployBucket(resourceId string) (nats.KeyValue, error) {
	return c.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: fmt.Sprintf("resources-%s-deploys", resourceId),
	})
}
