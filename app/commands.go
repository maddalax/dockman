package app

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/nats-io/nats.go"
)

type RunResourceCommand struct {
	ResourceId   string
	ResponseData *RunResourceResponse
	ResponseErr  error
}

type RunResourceResponse struct {
	Message string
	Error   error
}

func (c *RunResourceCommand) Execute(agent *Agent) {
	_, err := ResourceStart(agent.locator, c.ResourceId, StartOpts{
		RemoveExisting: true,
	})
	if err != nil {
		c.ResponseErr = err
	}
	fmt.Printf("Running resource: %s\n", c.ResourceId)
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
	_, err := ResourceStop(agent.locator, c.ResourceId)
	if err != nil {
		c.ResponseErr = err
	}
	fmt.Printf("Stopping resource: %s\n", c.ResourceId)
	c.ResponseData = &StopResourceResponse{
		Message: "Resource stopped",
	}
}

func (c *StopResourceCommand) GetResponse() any {
	return c.ResponseData
}

type StreamRunLogsCommand struct {
	ResourceId string
	SocketId   string
}

func (c *StreamRunLogsCommand) Execute(agent *Agent) {
	ctx := h.RequestContext{}
	ctx.Set("session-id", c.SocketId)
	context := WithAliveContext(&ctx)
	resource, err := ResourceGet(agent.locator, c.ResourceId)
	if err != nil {
		return
	}
	StreamResourceLogs(agent.locator, context, resource, func(msg *nats.Msg) {

	})
}

func (c *StreamRunLogsCommand) GetResponse() any {
	return nil
}
