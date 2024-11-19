package app

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/app/util/json2"
)

func ResourceList(locator *service.Locator) ([]*Resource, error) {
	client := service.Get[KvClient](locator)
	bucket, err := client.GetBucket("resources")
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
