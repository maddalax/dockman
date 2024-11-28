package app

import (
	"dockman/app/logger"
	"github.com/docker/docker/api/types"
)

type RunResourceCommand struct {
	ResourceId      string
	IgnoreIfRunning bool
	// if we change the instances and existing containers already exist for the new instance indexes, remove them
	RemoveExisting bool
	ResponseData   *RunResourceResponse
}

type RunResourceResponse struct {
	Message string
	Error   string
}

func (c *RunResourceCommand) Execute(agent *Agent) {
	_, err := ResourceStart(agent, c.ResourceId, StartOpts{
		RemoveExisting:  c.RemoveExisting,
		IgnoreIfRunning: c.IgnoreIfRunning,
	})
	if err != nil {
		c.ResponseData = &RunResourceResponse{
			Error: err.Error(),
		}
	} else {
		logger.InfoWithFields("Running resource", map[string]any{
			"resource_id": c.ResourceId,
		})
		c.ResponseData = &RunResourceResponse{
			Message: "Resource started",
		}
	}
}

func (c *RunResourceCommand) GetResponse() any {
	return c.ResponseData
}

func (c *RunResourceCommand) Name() string {
	return "RunResource"
}

type StopResourceCommand struct {
	ResourceId   string
	ResponseData *StopResourceResponse
}

type StopResourceResponse struct {
	Message string
	Error   string
}

func (c *StopResourceCommand) Execute(agent *Agent) {
	_, err := ResourceStop(agent, c.ResourceId)
	if err != nil {
		c.ResponseData = &StopResourceResponse{
			Error: err.Error(),
		}
	} else {
		logger.InfoWithFields("stopping resource", map[string]any{
			"resource_id": c.ResourceId,
		})
		c.ResponseData = &StopResourceResponse{
			Message: "Resource stopped",
		}
	}
}

func (c *StopResourceCommand) GetResponse() any {
	return c.ResponseData
}

func (c *StopResourceCommand) Name() string {
	return "StopResource"
}

type SetServerConfigCommand struct {
	Key   string
	Value string
}

func (c *SetServerConfigCommand) Execute(agent *Agent) {
	manager := agent.registry.GetServerConfigManager()
	manager.WriteConfig(c.Key, c.Value)
}

func (c *SetServerConfigCommand) GetResponse() any {
	return nil
}

func (c *SetServerConfigCommand) Name() string {
	return "SetServerConfig"
}

type GetServerConfigResponse struct {
	Value string
}

type GetServerConfigCommand struct {
	Key          string
	ResponseData GetServerConfigResponse
}

func (c *GetServerConfigCommand) Execute(agent *Agent) {
	manager := agent.registry.GetServerConfigManager()
	value := manager.GetConfig(c.Key)
	c.ResponseData = GetServerConfigResponse{
		Value: value,
	}
}

func (c *GetServerConfigCommand) GetResponse() any {
	return c.ResponseData
}

func (c *GetServerConfigCommand) Name() string {
	return "GetServerConfig"
}

type PingCommand struct {
	ResponseData *PingResponse
}

type PingResponse struct {
	Message string
}

func (p *PingCommand) Execute(agent *Agent) {
	p.ResponseData = &PingResponse{
		Message: "pong",
	}
}

func (p *PingCommand) GetResponse() any {
	return p.ResponseData
}

func (p *PingCommand) Name() string {
	return "Ping"
}

type GetContainerCommand struct {
	ResourceId   string
	Index        int
	ResponseData *GetContainerResponse
}

type GetContainerResponse struct {
	Error     error
	Container types.ContainerJSON
}

func (c *GetContainerCommand) Execute(agent *Agent) {
	resource, err := ResourceGet(agent.locator, c.ResourceId)
	if err != nil {
		c.ResponseData = &GetContainerResponse{
			Error: err,
		}
		return
	}
	switch resource.RunType {
	case RunTypeDockerBuild, RunTypeDockerRegistry:
		dockerClient, err := DockerConnect(agent.locator)
		if err != nil {
			c.ResponseData = &GetContainerResponse{
				Error: DockerConnectionError,
			}
			return
		}
		container, err := dockerClient.GetContainer(resource, c.Index)

		if err != nil {
			c.ResponseData = &GetContainerResponse{
				Error: err,
			}
			return
		}

		c.ResponseData = &GetContainerResponse{
			Container: container,
		}

	default:
		panic("Unsupported run type")
	}
}

func (c *GetContainerCommand) GetResponse() any {
	return c.ResponseData
}

func (c *GetContainerCommand) Name() string {
	return "GetContainer"
}
