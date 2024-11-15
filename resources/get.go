package resources

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"paas/domain"
	"paas/kv"
	"strings"
)

func Get(locator *service.Locator, id string) (*domain.Resource, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	client := service.Get[kv.Client](locator)

	if !strings.HasPrefix(id, "resources-") {
		id = fmt.Sprintf("resources-%s", id)
	}

	resourceBucket, err := client.GetBucket(id)

	if err != nil {
		return nil, err
	}

	return MapToResource(resourceBucket)
}
