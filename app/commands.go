package app

import (
	"dockside/app/logger"
)

type RunResourceCommand struct {
	ResourceId      string
	IgnoreIfRunning bool
	// if we change the instances and existing containers already exist for the new instance indexes, remove them
	RemoveExisting bool
	ResponseData   *RunResourceResponse
	ResponseErr    error
}

type RunResourceResponse struct {
	Message string
	Error   error
}

func (c *RunResourceCommand) Execute(agent *Agent) {
	_, err := ResourceStart(agent, c.ResourceId, StartOpts{
		RemoveExisting:  c.RemoveExisting,
		IgnoreIfRunning: c.IgnoreIfRunning,
	})
	if err != nil {
		c.ResponseErr = err
	}
	logger.InfoWithFields("Running resource", map[string]any{
		"resource_id": c.ResourceId,
	})
	c.ResponseData = &RunResourceResponse{
		Message: "Resource started",
	}
}

func (c *RunResourceCommand) GetResponse() any {
	return c.ResponseData
}

type StopResourceCommand struct {
	ResourceId   string
	ResponseData *StopResourceResponse
	ResponseErr  error
}

type StopResourceResponse struct {
	Message string
	Error   error
}

func (c *StopResourceCommand) Execute(agent *Agent) {
	_, err := ResourceStop(agent, c.ResourceId)
	if err != nil {
		c.ResponseErr = err
	}
	logger.InfoWithFields("stopping resource", map[string]any{
		"resource_id": c.ResourceId,
	})
	c.ResponseData = &StopResourceResponse{
		Message: "Resource stopped",
	}
}

func (c *StopResourceCommand) GetResponse() any {
	return c.ResponseData
}

type SetServerConfigCommand struct {
	Key   string
	Value string
}

func (c *SetServerConfigCommand) Execute(agent *Agent) {
	manager := agent.serverConfigManager
	manager.WriteConfig(c.Key, c.Value)
}

func (c *SetServerConfigCommand) GetResponse() any {
	return nil
}

type GetServerConfigResponse struct {
	Value string
}

type GetServerConfigCommand struct {
	Key          string
	ResponseData GetServerConfigResponse
}

func (c *GetServerConfigCommand) Execute(agent *Agent) {
	manager := agent.serverConfigManager
	value := manager.GetConfig(c.Key)
	c.ResponseData = GetServerConfigResponse{
		Value: value,
	}
}

func (c *GetServerConfigCommand) GetResponse() any {
	return c.ResponseData
}
