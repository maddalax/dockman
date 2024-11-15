package docker

import (
	"context"
	"fmt"
	"paas/domain"
)

func (c *Client) GetRunStatus(resource *domain.Resource) (domain.RunStatus, error) {
	containerName := fmt.Sprintf("%s-%s-container", resource.Name, resource.Id)
	inspect, err := c.cli.ContainerInspect(context.Background(), containerName)

	// unable to inspect, must not be running or docker is down
	if err != nil {
		return domain.RunStatusNotRunning, nil
	}

	if inspect.State != nil {
		if inspect.State.Running {
			return domain.RunStatusRunning, nil
		}
		if inspect.State.ExitCode != 0 {
			return domain.RunStatusErrored, nil
		}
		return domain.RunStatusNotRunning, nil
	}

	return domain.RunStatusUnknown, nil
}
