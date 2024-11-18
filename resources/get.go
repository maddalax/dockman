package resources

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"paas/domain"
	"paas/json2"
	"paas/kv"
	"strings"
)

func Get(locator *service.Locator, id string) (*domain.Resource, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	client := service.Get[kv.Client](locator)

	if strings.HasPrefix(id, "resources-") {
		id = strings.TrimPrefix(id, "resources-")
	}

	resourceBucket, err := client.GetBucket("resources")

	if err != nil {
		return nil, err
	}

	resource, err := resourceBucket.Get(id)

	if err != nil {
		return nil, err
	}

	return json2.Deserialize[domain.Resource](resource.Value())
}
