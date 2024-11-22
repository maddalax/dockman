package app

func (a *Agent) GetCurrentResourceServer(resource *Resource) *ResourceServer {
	for i, detail := range resource.ServerDetails {
		if detail.ServerId == a.serverId {
			return &resource.ServerDetails[i]
		}
	}
	return nil
}

func (a *Agent) CalculateResourceServer(resource *Resource) (*ResourceServer, error) {
	s := ResourceServer{}
	server, err := ServerGet(a.locator, a.serverId)
	if err != nil {
		return nil, err
	}
	switch resource.RunType {
	case RunTypeDockerBuild, RunTypeDockerRegistry:
		for i := range resource.InstancesPerServer {
			a.calculateDockerUpstreams(resource, server, &s, i)
		}
		s.RunStatus = a.GetRunStatus(resource)
	default:
		panic("unhandled default case")
	}
	return &s, nil
}

func (a *Agent) calculateDockerUpstreams(resource *Resource, server *Server, resourceServer *ResourceServer, index int) error {
	client, err := DockerConnect(a.locator)

	if err != nil {
		return err
	}

	container, err := client.GetContainer(resource, index)

	if err != nil {
		return err
	}

	hostIp := ""

	if server.RemoteIpAddress != "" {
		hostIp = server.RemoteIpAddress
	}

	// route using local ip first if possible
	if server.LocalIpAddress != "" {
		hostIp = server.LocalIpAddress
	}

	for port, binding := range container.NetworkSettings.Ports {
		if port.Proto() == "tcp" {
			for _, portBinding := range binding {
				upstream := Upstream{
					Host: hostIp,
					Port: portBinding.HostPort,
				}
				resourceServer.Upstreams = append(resourceServer.Upstreams, upstream)
			}
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
