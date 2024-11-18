package kv

import (
	"github.com/nats-io/nats.go"
)

func PutJson(bucket nats.KeyValue, key string, value interface{}) error {
	_, err := bucket.Put(key, MustSerialize(value))
	return err
}
