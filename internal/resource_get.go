package internal

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"paas/internal/util/json2"
	"strings"
)

func ResourceGet(locator *service.Locator, id string) (*Resource, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	client := service.Get[KvClient](locator)

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

	return json2.Deserialize[Resource](resource.Value())
}
