package resources

import (
	"errors"
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

	if patch["run_status"] != nil {
		resource.RunStatus = patch["run_status"].(domain.RunStatus)
	}

	if patch["instances_per_server"] != nil {
		resource.InstancesPerServer = patch["instances_per_server"].(int)
	}

	if patch["name"] != nil {
		return errors.New("name cannot be changed")
	}

	if patch["build_meta"] != nil {
		return errors.New("build meta cannot be changed")
	}

	validators := []Validator{
		RequiredFieldsValidator{
			resource: resource,
		},
	}

	for _, validator := range validators {
		err := validator.Validate()
		if err != nil {
			return err
		}
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
