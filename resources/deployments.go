package resources

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/kv"
	"time"
)

type CreateDeploymentRequest struct {
	ResourceId string
	BuildId    string
}

type Deployment struct {
	ResourceId string
	CreatedAt  time.Time
	BuildId    string
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

func CreateDeployment(locator *service.Locator, request CreateDeploymentRequest) error {
	client := service.Get[kv.Client](locator)

	bucket, _ := client.GetResourceDeployBucket(request.ResourceId)

	_, err := bucket.Put(request.BuildId, kv.MustSerialize(Deployment{
		ResourceId: request.ResourceId,
		CreatedAt:  time.Now(),
		BuildId:    request.BuildId,
	}))

	if err != nil {
		return err
	}

	return nil
}
