package docker

import (
	"context"
	"fmt"
	"paas/resources"
)

func (c *Client) GetRunStatus(resource *resources.Resource) (resources.RunStatus, error) {
	containerName := fmt.Sprintf("%s-%s-container", resource.Name, resource.Id)
	inspect, err := c.cli.ContainerInspect(context.Background(), containerName)

	// unable to inspect, must not be running or docker is down
	if err != nil {
		return resources.RunStatusNotRunning, nil
	}

	if inspect.State != nil {
		if inspect.State.Running {
			return resources.RunStatusRunning, nil
		}
		if inspect.State.ExitCode != 0 {
			return resources.RunStatusErrored, nil
		}
		return resources.RunStatusNotRunning, nil
	}

	return resources.RunStatusUnknown, nil
}
