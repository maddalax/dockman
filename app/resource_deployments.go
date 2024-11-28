package app

import (
	"dockman/app/util/json2"
	"encoding/json"
	"github.com/maddalax/htmgo/framework/service"
	"slices"
	"time"
)

func PatchDeployment(locator *service.Locator, resourceId string, buildId string, cb func(deployment *Deployment) *Deployment) error {
	client := service.Get[KvClient](locator)

	bucket, _ := client.GetResourceDeployBucket(resourceId)

	deployment, err := GetDeployment(locator, resourceId, buildId)

	if err != nil {
		return err
	}

	// apply the patch
	deployment = cb(deployment)

	_, err = bucket.Put(deployment.BuildId, json2.SerializeOrEmpty(deployment))

	if err != nil {
		return err
	}

	return nil
}

func GetDeployments(locator *service.Locator, resourceId string) ([]Deployment, error) {
	client := service.Get[KvClient](locator)
	bucket, err := client.GetResourceDeployBucket(resourceId)

	if err != nil {
		return nil, err
	}

	buildIds, err := bucket.ListKeys()

	if err != nil {
		return nil, err
	}

	var mapped []Deployment

	for buildId := range buildIds.Keys() {
		value, err := bucket.Get(buildId)
		if err != nil {
			continue
		}
		d := &Deployment{}
		err = json.Unmarshal(value.Value(), d)
		if err != nil {
			continue
		}
		mapped = append(mapped, *d)
	}

	slices.SortFunc(mapped, func(a, b Deployment) int {
		return int(b.CreatedAt.Sub(a.CreatedAt).Milliseconds())
	})

	return mapped, nil
}

func GetDeployment(locator *service.Locator, resourceId string, buildId string) (*Deployment, error) {
	client := service.Get[KvClient](locator)
	bucket, err := client.GetResourceDeployBucket(resourceId)

	if err != nil {
		return nil, err
	}

	build, err := bucket.Get(buildId)
	if err != nil {
		return nil, err
	}

	d := &Deployment{}
	err = json.Unmarshal(build.Value(), d)

	if err != nil {
		return nil, err
	}

	return d, nil
}

func CreateDeployment(locator *service.Locator, request CreateDeploymentRequest) error {
	client := service.Get[KvClient](locator)
	bucket, _ := client.GetResourceDeployBucket(request.ResourceId)

	_, err := bucket.Put(request.BuildId, json2.SerializeOrEmpty(Deployment{
		ResourceId: request.ResourceId,
		CreatedAt:  time.Now(),
		BuildId:    request.BuildId,
		Status:     DeploymentStatusPending,
		Source:     request.Source,
	}))

	if err != nil {
		return err
	}

	return nil
}
