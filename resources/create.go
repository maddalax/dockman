package resources

import (
	"github.com/google/uuid"
	"github.com/maddalax/htmgo/framework/service"
	"paas/history"
	"paas/kv/subject"
)

type CreateOptions struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	RunType     RunType           `json:"run_type"`
	BuildMeta   any               `json:"build_meta"`
	Env         map[string]string `json:"env"`
}

func Create(locator *service.Locator, options CreateOptions) (string, error) {
	resource := NewResource(uuid.NewString())

	resource.Name = options.Name
	resource.Environment = options.Environment
	resource.RunType = options.RunType
	resource.BuildMeta = options.BuildMeta
	resource.Env = options.Env

	validators := []Validator{
		BuildMetaValidator{
			meta: resource.BuildMeta,
		},
		RequiredFieldsValidator{
			resource: resource,
		},
	}

	for _, validator := range validators {
		err := validator.Validate()
		if err != nil {
			return "", err
		}
	}

	err := resource.Create(locator)

	if err != nil {
		return "", err
	}

	history.LogChange(locator, subject.ResourceCreated, map[string]any{
		"id":          resource.Id,
		"environment": options.Environment,
		"name":        options.Name,
	})

	return resource.Id, nil
}
