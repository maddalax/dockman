package resources

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/domain"
	"paas/kv"
	"time"
)

func UpdateDeploymentStatus(locator *service.Locator, request domain.UpdateDeploymentStatusRequest) error {
	deployment, err := GetDeployment(locator, request.ResourceId, request.BuildId)
	if err != nil {
		return err
	}
	deployment.Status = request.Status
	return SetDeployment(locator, *deployment)
}

func SetDeployment(locator *service.Locator, deployment domain.Deployment) error {
	client := service.Get[kv.Client](locator)

	bucket, _ := client.GetResourceDeployBucket(deployment.ResourceId)

	_, err := bucket.Put(deployment.BuildId, kv.MustSerialize(deployment))

	if err != nil {
		return err
	}

	return nil
}

func GetDeployments(locator *service.Locator, resourceId string) ([]domain.Deployment, error) {
	client := service.Get[kv.Client](locator)
	bucket, err := client.GetResourceDeployBucket(resourceId)

	if err != nil {
		return nil, err
	}

	buildIds, err := bucket.ListKeys()

	if err != nil {
		return nil, err
	}

	var mapped []domain.Deployment

	for buildId := range buildIds.Keys() {
		value, err := bucket.Get(buildId)
		if err != nil {
			continue
		}
		json := string(value.Value())
		d := kv.MustMapStringToStructure[domain.Deployment](json)
		if d == nil {
			continue
		}
		mapped = append(mapped, *d)
	}

	return mapped, nil
}

func GetDeployment(locator *service.Locator, resourceId string, buildId string) (*domain.Deployment, error) {
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
	d := kv.MustMapStringToStructure[domain.Deployment](json)

	return d, nil
}

func CreateDeployment(locator *service.Locator, request domain.CreateDeploymentRequest) error {
	return SetDeployment(locator, domain.Deployment{
		ResourceId: request.ResourceId,
		CreatedAt:  time.Now(),
		BuildId:    request.BuildId,
		Status:     domain.DeploymentStatusPending,
	})
}
