package resources

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/maddalax/htmgo/framework/service"
	"paas/history"
	"paas/kv"
	"paas/kv/subject"
)

type CreateOptions struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	RunType     RunType           `json:"run_type"`
	BuildMeta   any               `json:"build_meta"`
	Env         map[string]string `json:"env"`
}

func Create(locator *service.Locator, options CreateOptions) error {
	client := service.Get[kv.Client](locator)
	id := uuid.NewString()
	resourceBucketKey := fmt.Sprintf("resources-%s", id)

	bucket, err := client.GetBucket("resources")

	if err != nil {
		return err
	}

	_, err = bucket.Create(resourceBucketKey, []byte{})

	if err != nil {
		return err
	}

	resourceBucket, err := client.GetBucket(resourceBucketKey)

	if err != nil {
		return err
	}

	err = kv.AtomicPutMany(resourceBucket, func(m map[string]any) error {
		m["id"] = id
		m["environment"] = options.Environment
		m["run_type"] = options.RunType
		m["build_meta"] = options.BuildMeta
		m["name"] = options.Name
		for k, v := range options.Env {
			m[fmt.Sprintf("env/%s", k)] = v
		}
		return nil
	})

	if err != nil {
		return err
	}

	history.LogChange(locator, subject.ResourceCreated, map[string]any{
		"id":          id,
		"environment": options.Environment,
		"name":        options.Name,
	})

	return nil
}
