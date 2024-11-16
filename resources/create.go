package resources

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"paas/domain"
	"paas/history"
	"paas/kv"
	"paas/kv/subject"
)

type CreateOptions struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	RunType     domain.RunType    `json:"run_type"`
	BuildMeta   any               `json:"build_meta"`
	Env         map[string]string `json:"env"`
}

func Create(locator *service.Locator, options CreateOptions) (string, error) {
	client := service.Get[kv.Client](locator)
	resource := domain.NewResource(uuid.NewString())

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

	bucket, _ := client.GetBucket(resource.BucketKey())

	if bucket != nil {
		return "", fmt.Errorf("resource already exists")
	}

	bucket, err := client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: resource.BucketKey(),
	})

	if err != nil {
		return "", err
	}

	err = kv.AtomicPutMany(bucket, func(m map[string]any) error {
		m["id"] = resource.Id
		m["environment"] = resource.Environment
		m["run_type"] = resource.RunType
		m["build_meta"] = resource.BuildMeta
		m["name"] = resource.Name
		for k, v := range resource.Env {
			m[fmt.Sprintf("env/%s", k)] = v
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	bucket, err = client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "resources",
	})

	if err != nil {
		return "", err
	}

	// add it to the resources bucket for listing
	_, err = bucket.Create(resource.Id, []byte{})

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
