package internal

import (
	"fmt"
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
