package docker

import (
	"paas/domain"
)

func (c *Client) GetRunStatus(resource *domain.Resource) (domain.RunStatus, error) {
	statuses := make([]domain.RunStatus, resource.InstancesPerServer)

	for i := range resource.InstancesPerServer {
		inspect, err := c.GetContainer(resource, i)
		if err != nil {
			statuses[i] = domain.RunStatusNotRunning
			continue
		}

		if inspect.State != nil {
			if inspect.State.Running {
				statuses[i] = domain.RunStatusRunning
			} else {
				statuses[i] = domain.RunStatusNotRunning
			}
		} else {
			statuses[i] = domain.RunStatusUnknown
		}
	}

	allRunning := true
	anyRunning := false
	for _, status := range statuses {
		if status != domain.RunStatusRunning {
			allRunning = false
		}
		if status == domain.RunStatusRunning {
			anyRunning = true
		}
	}

	if allRunning {
		return domain.RunStatusRunning, nil
	}

	if anyRunning {
		return domain.RunStatusPartiallyRunning, nil
	}

	return domain.RunStatusNotRunning, nil
}
