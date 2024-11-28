package app

import (
	"context"
	"dockside/app/logger"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"strconv"
	"strings"
)

type RunningContainer struct {
	Container   types.Container
	Index       int
	NameNoIndex string
	ResourceId  string
}

func (c *DockerClient) GetRunningContainers(resource *Resource) ([]RunningContainer, error) {
	containers, err := c.cli.ContainerList(context.Background(), container.ListOptions{})

	if err != nil {
		return nil, err
	}

	var matched []RunningContainer
	containerNameNoIndex := fmt.Sprintf("/%s-%s-container-", resource.Name, resource.Id)

	for _, t := range containers {
		if len(t.Names) == 0 {
			continue
		}
		name := t.Names[0]
		if strings.HasPrefix(name, containerNameNoIndex) {
			containerIndex, err := strconv.Atoi(t.Names[0][len(containerNameNoIndex):])
			if err != nil {
				logger.InfoWithFields("Error parsing container index", map[string]any{
					"container_id":   t.ID,
					"resource_id":    resource.Id,
					"container_name": name,
				})
				continue
			}
			matched = append(matched, RunningContainer{
				Container:   t,
				Index:       containerIndex,
				NameNoIndex: containerNameNoIndex,
				ResourceId:  resource.Id,
			})
		}
	}

	return matched, nil
}

// ReduceToMatchResourceCount reduces the number of containers to match the resource count
func (c *DockerClient) ReduceToMatchResourceCount(resource *Resource, count int) {
	containers, err := c.GetRunningContainers(resource)

	if err != nil {
		logger.ErrorWithFields("Error getting running containers", err, map[string]any{
			"resource_id": resource.Id,
		})
		return
	}

	for _, t := range containers {
		if err != nil {
			continue
		}
		if t.Index >= count {
			logger.InfoWithFields("Stopping container", map[string]any{
				"container_id": t.Container.ID,
				"name":         t.Container.Names[0],
			})
			err = c.cli.ContainerStop(context.Background(), t.Container.ID, container.StopOptions{})
			if err != nil {
				logger.ErrorWithFields("Error stopping container", err, map[string]any{
					"container_id": t.Container.ID,
				})
			}
			err = c.cli.ContainerRemove(context.Background(), t.Container.ID, container.RemoveOptions{})
			if err != nil {
				logger.ErrorWithFields("Error removing container", err, map[string]any{
					"container_id": t.Container.ID,
				})
			}
		}
	}
}
