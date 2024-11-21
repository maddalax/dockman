package app

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"io"
)

type ImageStore struct {
	store nats.ObjectStore
}

func NewImageStore(store nats.ObjectStore) *ImageStore {
	return &ImageStore{
		store: store,
	}
}

func (s *ImageStore) Get(imageId string) (nats.ObjectResult, error) {
	return s.store.Get(imageId)
}

func (s *ImageStore) Has(imageId string) bool {
	_, err := s.store.GetInfo(imageId)
	return err == nil
}

func (s *ImageStore) ImageIdForResource(resource *Resource) string {
	return fmt.Sprintf("%s-%s", resource.Name, resource.Id)
}

func (s *ImageStore) GetBuildId(imageId string) string {
	info, err := s.store.GetInfo(imageId)

	if err != nil {
		return ""
	}

	return info.Metadata["buildId"]
}

func (s *ImageStore) Put(obj *nats.ObjectMeta, reader io.Reader, opts ...nats.ObjectOpt) (*nats.ObjectInfo, error) {
	return s.store.Put(obj, reader, opts...)
}

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

func (c *KvClient) ImageStore() (*ImageStore, error) {
	bucket, err := c.GetOrCreateObjectStore(&nats.ObjectStoreConfig{
		Bucket: "images",
	})
	if err != nil {
		return nil, err
	}
	return NewImageStore(bucket), nil
}
