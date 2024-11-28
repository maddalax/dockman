package app

import (
	"dockman/app/util/json2"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
)

func ResourceList(locator *service.Locator) ([]*Resource, error) {
	client := service.Get[KvClient](locator)
	bucket, err := client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "resources",
	})
	resources := make([]*Resource, 0)

	if err != nil {
		return nil, err
	}

	listener, err := bucket.ListKeys()

	if err != nil {
		return nil, err
	}

	for s := range listener.Keys() {
		resource, err := bucket.Get(s)
		if err != nil {
			continue
		}
		mapped, err := json2.Deserialize[Resource](resource.Value())
		if err != nil {
			continue
		}
		resources = append(resources, mapped)
	}
	return resources, nil
}
