package resources

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"paas/kv"
)

type Resource struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	RunType     RunType           `json:"run_type"`
	BuildMeta   any               `json:"build_meta"`
	Env         map[string]string `json:"env"`
	RunStatus   RunStatus         `json:"run_status"`
}

type RunStatus = int

const (
	RunStatusUnknown RunStatus = iota
	RunStatusNotRunning
	RunStatusRunning
	RunStatusErrored
)

func NewResource(id string) *Resource {
	resource := &Resource{
		Id: id,
	}
	return resource
}

func (resource *Resource) GetBucketKey() string {
	return fmt.Sprintf("resources-%s", resource.Id)
}

func (resource *Resource) GetBucket(locator *service.Locator) (nats.KeyValue, error) {
	client := service.Get[kv.Client](locator)
	return client.GetBucket(resource.GetBucketKey())
}

func (resource *Resource) Create(locator *service.Locator) error {
	bucket, _ := resource.GetBucket(locator)
	if bucket != nil {
		return fmt.Errorf("resource already exists")
	}
	client := service.Get[kv.Client](locator)
	key := resource.GetBucketKey()
	bucket, err := client.GetOrCreateBucket(key)
	if err != nil {
		return err
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
		return err
	}

	// add it to the resources bucket for listing
	bucket, err = client.GetOrCreateBucket("resources")

	if err != nil {
		return err
	}

	_, err = bucket.Create(key, []byte{})

	return nil
}

func (resource *Resource) SetRunStatus(locator *service.Locator, status RunStatus) error {
	return resource.Patch(locator, map[string]any{
		"run_status": status,
	})
}

func (resource *Resource) Patch(locator *service.Locator, data map[string]any) error {
	bucket, err := resource.GetBucket(locator)
	if err != nil {
		return err
	}
	err = kv.AtomicPutMany(bucket, func(m map[string]any) error {
		for k, v := range data {
			m[k] = v
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
