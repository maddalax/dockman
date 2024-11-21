package app

func (a *Agent) GetResourceServer(resource *Resource) *ResourceServer {
	for i, detail := range resource.ServerDetails {
		if detail.ServerId == a.serverId {
			return &resource.ServerDetails[i]
		}
	}
	return nil
}

func (a *Agent) GetRunStatus(resource *Resource) RunStatus {
	if resource.RunType == RunTypeDockerBuild || resource.RunType == RunTypeDockerRegistry {
		return a.getRunStatusDocker(resource)
	}
	return RunStatusUnknown
}

func (a *Agent) getRunStatusDocker(resource *Resource) RunStatus {
	client, err := DockerConnect(a.locator)
	if err != nil {
		return RunStatusNotRunning
	}
	status, err := client.GetRunStatus(resource)
	if err != nil {
		return RunStatusNotRunning
	}
	return status
}
