package resources

import (
	"github.com/maddalax/htmgo/framework/service"
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
