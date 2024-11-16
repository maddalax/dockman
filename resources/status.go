package resources

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"paas/docker"
	"paas/domain"
	"paas/history"
	"paas/kv/subject"
	"paas/util"
	"time"
)

// IsRunnable checks if a resource is runnable
// in the case of docker, the container must exist
func IsRunnable(resource *domain.Resource) bool {
	switch resource.RunType {
	case domain.RunTypeDockerBuild:
		fallthrough
	case domain.RunTypeDockerRegistry:
		client, err := docker.Connect()
		if err != nil {
			return false
		}
		_, err = client.GetContainer(resource)
		return err == nil
	default:
		return false
	}
}

func Start(locator *service.Locator, resourceId string) (*domain.Resource, error) {
	lock := GetStatusLock(locator, resourceId)
	err := lock.Lock()
	if err != nil {
		return nil, err
	}
	defer lock.Unlock()

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
			RemoveExisting: false,
		})
		if err != nil {
			return nil, err
		}
		success := waitForStatus(locator, resourceId, domain.RunStatusRunning)
		if !success {
			return nil, domain.ResourceFailedToStartError
		}
		resource, err = Get(locator, resourceId)
		if err != nil {
			return nil, err
		}
		return resource, nil
	default:
		return nil, domain.UnsupportedRunTypeError
	}
}

func Stop(locator *service.Locator, resourceId string) (*domain.Resource, error) {
	lock := GetStatusLock(locator, resourceId)
	err := lock.Lock()
	if err != nil {
		return nil, err
	}
	defer lock.Unlock()

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
		success := waitForStatus(locator, resourceId, domain.RunStatusNotRunning)
		if !success {
			return nil, domain.ResourceFailedToStopError
		}
		resource, err = Get(locator, resourceId)
		if err != nil {
			return nil, err
		}
		return resource, nil
	default:
		return nil, domain.UnsupportedRunTypeError
	}
}

func waitForStatus(locator *service.Locator, resourceId string, status domain.RunStatus) bool {
	success := util.WaitFor(time.Second*5, 200*time.Millisecond, func() bool {
		resource, err := Get(locator, resourceId)
		if err != nil {
			return false
		}
		fmt.Printf("waiting for status %v, got %v\n", status, resource.RunStatus)
		return resource.RunStatus == status
	})
	return success
}
