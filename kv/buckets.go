package kv

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

func (c *Client) GetResourceDeployBucket(resourceId string) (nats.KeyValue, error) {
	return c.GetBucket(fmt.Sprintf("resources-%s-deploys", resourceId))
}
