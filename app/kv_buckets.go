package app

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"paas/app/util/must"
)

func (c *KvClient) GetResourceDeployBucket(resourceId string) (nats.KeyValue, error) {
	return c.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: fmt.Sprintf("resources-%s-deploys", resourceId),
	})
}

func (c *KvClient) PutJson(bucket nats.KeyValue, key string, value interface{}) error {
	_, err := bucket.Put(key, must.Serialize(value))
	return err
}
