package app

import (
	"errors"
	"github.com/nats-io/nats.go"
)

func (c *KvClient) GetOrCreateObjectStore(config *nats.ObjectStoreConfig) (nats.ObjectStore, error) {
	store, err := c.js.ObjectStore(config.Bucket)
	if err != nil {
		if errors.Is(err, nats.ErrStreamNotFound) {
			store, err = c.js.CreateObjectStore(config)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return store, nil
}
