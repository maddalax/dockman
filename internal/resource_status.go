package internal

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"paas/internal/subject"
	"paas/internal/util"
	"time"
)

// IsResourceRunnable checks if a resource is runnable
// in the case of docker, the container must exist for each instance
func IsResourceRunnable(resource *Resource) (bool, error) {
	switch resource.RunType {
	case RunTypeDockerBuild:
		fallthrough
	case RunTypeDockerRegistry:
		client, err := DockerConnect()
		if err != nil {
			return false, err
		}
		for i := range resource.InstancesPerServer {
			_, err = client.GetContainer(resource, i)
			if err != nil {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, nil
	}
}

type StartOpts struct {
	// Whether to ignore if the resource is already running
	IgnoreIfRunning bool
	// Whether to remove the existing instances before running a new one with the same id
	RemoveExisting bool
}

func ResourceStart(locator *service.Locator, resourceId string, opts StartOpts) (*Resource, error) {
	lock := ResourceStatusLock(locator, resourceId)
	err := lock.Lock()
	if err != nil {
		return nil, err
	}
	defer lock.Unlock()

	resource, err := ResourceGet(locator, resourceId)
	if err != nil {
		return nil, err
	}
	LogChange(locator, subject.ResourceStarted, map[string]any{
		"resource_id": resource.Id,
	})
	switch resource.RunType {
	case RunTypeDockerBuild:
		fallthrough
	case RunTypeDockerRegistry:
		client, err := DockerConnect()
		if err != nil {
			return nil, err
		}
		err = client.Run(resource, RunOptions{
			RemoveExisting:  opts.RemoveExisting,
			IgnoreIfRunning: opts.IgnoreIfRunning,
		})
		if err != nil {
			return nil, err
		}
		success := waitForStatus(locator, resourceId, RunStatusRunning)
		if !success {
			return nil, ResourceFailedToStartError
		}
		resource, err = ResourceGet(locator, resourceId)
		if err != nil {
			return nil, err
		}
		return resource, nil
	default:
		return nil, UnsupportedRunTypeError
	}
}

func ResourceStop(locator *service.Locator, resourceId string) (*Resource, error) {
	lock := ResourceStatusLock(locator, resourceId)
	err := lock.Lock()
	if err != nil {
		return nil, err
	}
	defer lock.Unlock()

	resource, err := ResourceGet(locator, resourceId)
	if err != nil {
		return nil, err
	}
	LogChange(locator, subject.ResourceStopped, map[string]any{
		"resource_id": resource.Id,
	})
	switch resource.RunType {
	case RunTypeDockerBuild:
		fallthrough
	case RunTypeDockerRegistry:
		client, err := DockerConnect()
		if err != nil {
			return nil, err
		}
		err = client.Stop(resource)
		if err != nil {
			return nil, err
		}
		success := waitForStatus(locator, resourceId, RunStatusNotRunning)
		if !success {
			return nil, ResourceFailedToStopError
		}
		resource, err = ResourceGet(locator, resourceId)
		if err != nil {
			return nil, err
		}
		return resource, nil
	default:
		return nil, UnsupportedRunTypeError
	}
}

func waitForStatus(locator *service.Locator, resourceId string, status RunStatus) bool {
	success := util.WaitFor(time.Second*5, 200*time.Millisecond, func() bool {
		resource, err := ResourceGet(locator, resourceId)
		if err != nil {
			return false
		}
		fmt.Printf("waiting for status %v, got %v\n", status, resource.RunStatus)
		return resource.RunStatus == status
	})
	return success
}
