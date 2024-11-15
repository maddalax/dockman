package resources

import (
	"errors"
	"github.com/maddalax/htmgo/framework/service"
	"paas/docker"
	"paas/domain"
	"paas/history"
	"paas/kv/subject"
)

func Start(locator *service.Locator, resourceId string) (*domain.Resource, error) {
	resource, err := Get(locator, resourceId)
	if err != nil {
		return nil, err
	}
	history.LogChange(locator, subject.ResourceStarted, map[string]any{
		"resource_id": resource.Id,
	})
	switch resource.RunType {
	case domain.RunTypeDockerBuild:
		fallthrough
	case domain.RunTypeDockerRegistry:
		client, err := docker.Connect()
		if err != nil {
			return nil, err
		}
		err = client.Run(resource, docker.RunOptions{
			KillExisting: true,
		})
		if err != nil {
			return nil, err
		}
		return resource, nil
	default:
		return nil, errors.New("unsupported run type")
	}
}

func Stop(locator *service.Locator, resourceId string) (*domain.Resource, error) {
	resource, err := Get(locator, resourceId)
	if err != nil {
		return nil, err
	}
	history.LogChange(locator, subject.ResourceStopped, map[string]any{
		"resource_id": resource.Id,
	})
	switch resource.RunType {
	case domain.RunTypeDockerBuild:
		fallthrough
	case domain.RunTypeDockerRegistry:
		client, err := docker.Connect()
		if err != nil {
			return nil, err
		}
		err = client.Stop(resource)
		if err != nil {
			return nil, err
		}
		return resource, nil
	default:
		return nil, errors.New("unsupported run type")
	}
}
