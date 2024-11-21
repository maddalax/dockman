package app

import (
	"dockside/app/subject"
	"github.com/google/uuid"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
)

type ResourceCreateOptions struct {
	Name               string            `json:"name"`
	Environment        string            `json:"environment"`
	RunType            RunType           `json:"run_type"`
	BuildMeta          BuildMeta         `json:"build_meta"`
	Env                map[string]string `json:"env"`
	InstancesPerServer int               `json:"instances_per_server"`
}

func ResourceCreate(locator *service.Locator, options ResourceCreateOptions) (string, error) {
	client := service.Get[KvClient](locator)
	resource := NewResource(uuid.NewString())

	resource.Name = options.Name
	resource.Environment = options.Environment
	resource.RunType = options.RunType
	resource.BuildMeta = options.BuildMeta
	resource.Env = options.Env
	resource.InstancesPerServer = options.InstancesPerServer

	if resource.InstancesPerServer == 0 {
		resource.InstancesPerServer = 1
	}

	validators := []Validator{
		BuildMetaValidator{
			Meta: resource.BuildMeta,
		},
		RequiredFieldsValidator{
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

	err = client.PutJson(bucket, resource.Id, resource)

	if err != nil {
		return "", err
	}

	LogChange(locator, subject.ResourceCreated, map[string]any{
		"id":          resource.Id,
		"environment": options.Environment,
		"name":        options.Name,
	})

	return resource.Id, nil
}
