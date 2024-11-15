package resources

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/domain"
	"paas/history"
	"paas/kv"
	"paas/kv/subject"
)

func SetRunStatus(locator *service.Locator, resourceId string, status domain.RunStatus) error {
	return Patch(locator, resourceId, map[string]any{
		"run_status": status,
	})
}

func Patch(locator *service.Locator, id string, patch map[string]any) error {
	resource, err := Get(locator, id)
	if err != nil {
		return err
	}

	client := service.Get[kv.Client](locator)
	bucket, err := client.GetBucket(resource.BucketKey())

	if err != nil {
		return err
	}

	err = kv.AtomicPutMany(bucket, func(m map[string]any) error {
		for k, v := range patch {
			m[k] = v
		}
		return nil
	})

	if err != nil {
		return err
	}

	patch["resource_id"] = resource.Id
	history.LogChange(locator, subject.ResourcePatched, patch)

	return nil
}
