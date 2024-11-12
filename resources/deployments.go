package resources

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/kv"
	"time"
)

type DeploymentStatus string

const (
	DeploymentStatusPending   DeploymentStatus = "pending"
	DeploymentStatusRunning   DeploymentStatus = "running"
	DeploymentStatusSucceeded DeploymentStatus = "succeeded"
	DeploymentStatusFailed    DeploymentStatus = "failed"
)

type CreateDeploymentRequest struct {
	ResourceId string
	BuildId    string
}

type UpdateDeploymentStatusRequest struct {
	ResourceId string
	BuildId    string
	Status     DeploymentStatus
}

type Deployment struct {
	ResourceId string
	CreatedAt  time.Time
	BuildId    string
	Status     DeploymentStatus
}

func UpdateDeploymentStatus(locator *service.Locator, request UpdateDeploymentStatusRequest) error {
	deployment, err := GetDeployment(locator, request.ResourceId, request.BuildId)
	if err != nil {
		return err
	}
	deployment.Status = request.Status
	return SetDeployment(locator, *deployment)
}

func SetDeployment(locator *service.Locator, deployment Deployment) error {
	client := service.Get[kv.Client](locator)

	bucket, _ := client.GetResourceDeployBucket(deployment.ResourceId)

	_, err := bucket.Put(deployment.BuildId, kv.MustSerialize(deployment))

	if err != nil {
		return err
	}

	return nil
}

func GetDeployments(locator *service.Locator, resourceId string) ([]Deployment, error) {
	client := service.Get[kv.Client](locator)
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
		json := string(value.Value())
		d := kv.MustMapStringToStructure[Deployment](json)
		if d == nil {
			continue
		}
		mapped = append(mapped, *d)
	}

	return mapped, nil
}

func GetDeployment(locator *service.Locator, resourceId string, buildId string) (*Deployment, error) {
	client := service.Get[kv.Client](locator)
	bucket, err := client.GetResourceDeployBucket(resourceId)

	if err != nil {
		return nil, err
	}

	build, err := bucket.Get(buildId)
	if err != nil {
		return nil, err
	}

	json := string(build.Value())
	d := kv.MustMapStringToStructure[Deployment](json)

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
