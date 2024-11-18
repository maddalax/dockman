package app

import (
	"github.com/nats-io/nats.go"
	"paas/must"
)

func PutJson(bucket nats.KeyValue, key string, value interface{}) error {
	_, err := bucket.Put(key, must.Serialize(value))
	return err
}
