package resources

import (
	"errors"
	"github.com/maddalax/htmgo/framework/service"
	"paas/domain"
	"paas/json2"
	"paas/kv"
	"paas/validation"
)

func SetRunStatus(locator *service.Locator, resourceId string, status domain.RunStatus) error {
	return Patch(locator, resourceId, func(resource *domain.Resource) *domain.Resource {
		resource.RunStatus = status
		return resource
	})
}

func Patch(locator *service.Locator, id string, cb func(resource *domain.Resource) *domain.Resource) error {
	lock := GetPatchLock(locator, id)
	err := lock.Lock()

	if err != nil {
		return err
	}

	defer lock.Unlock()

	resource, err := Get(locator, id)
	if err != nil {
		return err
	}

	updated := cb(resource)
	current, err := Get(locator, id)

	if err != nil {
		return err
	}

	validators := []validation.Validator{
		validation.RequiredFieldsValidator{
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

	_, err = current.BuildMeta.ValidatePatch(updated.BuildMeta)

	if err != nil {
		return err
	}

	client := kv.GetClientFromLocator(locator)

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
