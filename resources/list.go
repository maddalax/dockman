package resources

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/domain"
	"paas/json2"
	"paas/kv"
)

func List(locator *service.Locator) ([]*domain.Resource, error) {
	client := service.Get[kv.Client](locator)
	bucket, err := client.GetBucket("resources")
	resources := make([]*domain.Resource, 0)

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
		mapped, err := json2.Deserialize[domain.Resource](resource.Value())
		if err != nil {
			continue
		}
		resources = append(resources, mapped)
	}
	return resources, nil
}
