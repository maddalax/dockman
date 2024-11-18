package resources

import (
	"github.com/google/uuid"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"paas/domain"
	"paas/history"
	"paas/kv"
	"paas/kv/subject"
	"paas/validation"
)

type CreateOptions struct {
	Name               string            `json:"name"`
	Environment        string            `json:"environment"`
	RunType            domain.RunType    `json:"run_type"`
	BuildMeta          domain.BuildMeta  `json:"build_meta"`
	Env                map[string]string `json:"env"`
	InstancesPerServer int               `json:"instances_per_server"`
}

func Create(locator *service.Locator, options CreateOptions) (string, error) {
	client := service.Get[kv.Client](locator)
	resource := domain.NewResource(uuid.NewString())

	resource.Name = options.Name
	resource.Environment = options.Environment
	resource.RunType = options.RunType
	resource.BuildMeta = options.BuildMeta
	resource.Env = options.Env
	resource.InstancesPerServer = options.InstancesPerServer

	if resource.InstancesPerServer == 0 {
		resource.InstancesPerServer = 1
	}

	validators := []validation.Validator{
		validation.BuildMetaValidator{
			Meta: resource.BuildMeta,
		},
		validation.RequiredFieldsValidator{
			Resource: resource,
		},
	}

	for _, validator := range validators {
		err := validator.Validate()
		if err != nil {
			return "", err
		}
	}

	bucket, err := client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "resources",
	})

	err = kv.PutJson(bucket, resource.Id, resource)

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
