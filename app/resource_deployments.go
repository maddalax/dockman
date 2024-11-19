package app

import (
	"encoding/json"
	"github.com/maddalax/htmgo/framework/service"
	"paas/app/util/must"
	"time"
)

func UpdateDeploymentStatus(locator *service.Locator, request UpdateDeploymentStatusRequest) error {
	deployment, err := GetDeployment(locator, request.ResourceId, request.BuildId)
	if err != nil {
		return err
	}
	deployment.Status = request.Status
	return SetDeployment(locator, *deployment)
}

func SetDeployment(locator *service.Locator, deployment Deployment) error {
	client := service.Get[KvClient](locator)

	bucket, _ := client.GetResourceDeployBucket(deployment.ResourceId)

	_, err := bucket.Put(deployment.BuildId, must.Serialize(deployment))

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
	return SetDeployment(locator, Deployment{
		ResourceId: request.ResourceId,
		CreatedAt:  time.Now(),
		BuildId:    request.BuildId,
		Status:     DeploymentStatusPending,
	})
}
