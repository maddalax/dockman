package app

func (c *DockerClient) GetRunStatus(resource *Resource) (RunStatus, error) {
	statuses := make([]RunStatus, resource.InstancesPerServer)

	for i := range resource.InstancesPerServer {
		inspect, err := c.GetContainer(resource, i)
		if err != nil {
			statuses[i] = RunStatusNotRunning
			continue
		}

		if inspect.State != nil {
			if inspect.State.Running {
				statuses[i] = RunStatusRunning
			} else {
				statuses[i] = RunStatusNotRunning
			}
		} else {
			statuses[i] = RunStatusUnknown
		}
	}

	allRunning := true
	anyRunning := false
	for _, status := range statuses {
		if status != RunStatusRunning {
			allRunning = false
		}
		if status == RunStatusRunning {
			anyRunning = true
		}
	}

	if allRunning {
		return RunStatusRunning, nil
	}

	if anyRunning {
		return RunStatusPartiallyRunning, nil
	}

	return RunStatusNotRunning, nil
}
