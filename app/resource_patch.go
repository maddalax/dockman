package app

import (
	"errors"
	"github.com/maddalax/htmgo/framework/service"
	"paas/json2"
)

func SetRunStatus(locator *service.Locator, resourceId string, status RunStatus) error {
	return Patch(locator, resourceId, func(resource *Resource) *Resource {
		resource.RunStatus = status
		return resource
	})
}

func Patch(locator *service.Locator, id string, cb func(resource *Resource) *Resource) error {
	lock := GetPatchLock(locator, id)
	err := lock.Lock()

	if err != nil {
		return err
	}

	defer lock.Unlock()

	resource, err := ResourceGet(locator, id)
	if err != nil {
		return err
	}

	updated := cb(resource)
	current, err := ResourceGet(locator, id)

	if err != nil {
		return err
	}

	validators := []Validator{
		RequiredFieldsValidator{
			Resource: resource,
		},
	}

	for _, validator := range validators {
		err := validator.Validate()
		if err != nil {
			return err
		}
	}

	if updated.Name != current.Name {
		return errors.New("name cannot be changed")
	}

	err = current.BuildMeta.ValidatePatch(updated.BuildMeta)

	if err != nil {
		return err
	}

	client := GetClientFromLocator(locator)

	resources, err := client.GetBucket("resources")

	if err != nil {
		return err
	}

	serialized, err := json2.Serialize(updated)

	if err != nil {
		return err
	}

	_, err = resources.Put(id, serialized)

	return err
}
