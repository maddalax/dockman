package app

import (
	"dockside/app/logger"
	"dockside/app/subject"
	"dockside/app/util"
	"github.com/maddalax/htmgo/framework/service"
	"time"
)

// IsResourceRunnable checks if a resource is runnable
func IsResourceRunnable(locator *service.Locator, resource *Resource) (bool, error) {
	store, err := KvFromLocator(locator).ImageStore()

	if err != nil {
		return false, err
	}

	// if we have a built image in the store, then its runnable
	has := store.Has(store.ImageIdForResource(resource))

	return has, nil
}

type StartOpts struct {
	// Whether to ignore if the resource is already running
	IgnoreIfRunning bool
	// Whether to remove the existing instances before running a new one with the same id
	RemoveExisting bool
}

func SendResourceStartCommand(locator *service.Locator, resourceId string, opts StartOpts) ([]*SendCommandResponse[RunResourceResponse], error) {
	responses, err := SendCommandForResource[RunResourceResponse](locator, resourceId, SendCommandOpts{
		Command: &RunResourceCommand{
			ResourceId:      resourceId,
			IgnoreIfRunning: opts.IgnoreIfRunning,
			RemoveExisting:  opts.RemoveExisting,
		},
		// May take a while to start if it's a large container that needs to be downloaded
		Timeout: time.Second * 30,
	})
	return responses, err
}

func SendResourceStopCommand(locator *service.Locator, resourceId string) ([]*SendCommandResponse[StopResourceResponse], error) {
	responses, err := SendCommandForResource[StopResourceResponse](locator, resourceId, SendCommandOpts{
		Command: &StopResourceCommand{
			ResourceId: resourceId,
		},
		Timeout: time.Second * 10,
	})
	return responses, err
}

// ResourceStart starts a resource, blocking until the resource is started.
// Note: this should only be called from a command so it is propagated to all servers
func ResourceStart(agent *Agent, resourceId string, opts StartOpts) (*Resource, error) {
	lock := ResourceStatusLock(agent.locator, resourceId)
	err := lock.Lock()
	if err != nil {
		return nil, err
	}
	defer lock.Unlock()

	resource, err := ResourceGet(agent.locator, resourceId)
	if err != nil {
		return nil, err
	}
	LogChange(agent.locator, subject.ResourceStarted, map[string]any{
		"resource_id": resource.Id,
	})
	switch resource.RunType {
	case RunTypeDockerBuild:
		fallthrough
	case RunTypeDockerRegistry:
		client, err := DockerConnect(agent.locator)
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
		success := waitForStatus(agent, resourceId, RunStatusRunning)
		if !success {
			return nil, ResourceFailedToStartError
		}
		resource, err = ResourceGet(agent.locator, resourceId)
		if err != nil {
			return nil, err
		}
		return resource, nil
	default:
		return nil, UnsupportedRunTypeError
	}
}

// ResourceStop stops a resource, blocking until the resource is stopped.
// Note: this should only be called from a command so it is propagated to all servers
func ResourceStop(agent *Agent, resourceId string) (*Resource, error) {
	lock := ResourceStatusLock(agent.locator, resourceId)
	err := lock.Lock()
	if err != nil {
		return nil, err
	}
	defer lock.Unlock()

	resource, err := ResourceGet(agent.locator, resourceId)
	if err != nil {
		return nil, err
	}
	LogChange(agent.locator, subject.ResourceStopped, map[string]any{
		"resource_id": resource.Id,
	})
	switch resource.RunType {
	case RunTypeDockerBuild:
		fallthrough
	case RunTypeDockerRegistry:
		client, err := DockerConnect(agent.locator)
		if err != nil {
			return nil, err
		}
		err = client.Stop(resource)
		if err != nil {
			return nil, err
		}
		success := waitForStatus(agent, resourceId, RunStatusNotRunning)
		if !success {
			return nil, ResourceFailedToStopError
		}
		resource, err = ResourceGet(agent.locator, resourceId)
		if err != nil {
			return nil, err
		}
		return resource, nil
	default:
		return nil, UnsupportedRunTypeError
	}
}

func waitForStatus(agent *Agent, resourceId string, status RunStatus) bool {
	success := util.WaitFor(time.Second*10, time.Second, func() bool {
		resource, err := ResourceGet(agent.locator, resourceId)
		if err != nil {
			return false
		}
		logger.DebugWithFields("waiting for status", map[string]any{
			"status":      status,
			"resource_id": resourceId,
		})
		return agent.GetRunStatus(resource) == status
	})
	return success
}

func GetComputedRunStatus(resource *Resource) RunStatus {
	allRunning := true
	anyRunning := false

	for _, s := range resource.ServerDetails {
		if s.RunStatus != RunStatusRunning {
			allRunning = false
		}
		if s.RunStatus == RunStatusRunning || s.RunStatus == RunStatusPartiallyRunning {
			anyRunning = true
		}
	}

	if allRunning {
		return RunStatusRunning
	}
	if anyRunning {
		return RunStatusPartiallyRunning
	}
	return RunStatusNotRunning
}
