package resources

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/domain"
	"paas/kv"
)

type ResourceName struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func GetNames(locator *service.Locator) []ResourceName {
	client := service.Get[kv.Client](locator)
	bucket, err := client.GetBucket("resources")

	if err != nil {
		return []ResourceName{}
	}

	listener, err := bucket.ListKeys()

	if err != nil {
		return []ResourceName{}
	}

	mapped, err := kv.MustMapIntoMany[ResourceName](client, listener.Keys(), "name", "id")

	if err != nil {
		return []ResourceName{}
	}

	return mapped

}

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
		resourceBucket, err := client.GetBucket(s)
		if err != nil {
			continue
		}
		mapped, err := MapToResource(resourceBucket)
		if err != nil {
			continue
		}
		resources = append(resources, mapped)
	}
	return resources, nil
}
