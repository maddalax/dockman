package app

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"log/slog"
	"strconv"
	"strings"
)

// ReduceToMatchResourceCount reduces the number of containers to match the resource count
func (c *DockerClient) ReduceToMatchResourceCount(resource *Resource, count int) {
	containers, err := c.cli.ContainerList(context.Background(), container.ListOptions{})

	if err != nil {
		return
	}

	var matched []*types.Container
	containerNameNoIndex := fmt.Sprintf("/%s-%s-container-", resource.Name, resource.Id)

	for _, t := range containers {
		if len(t.Names) == 0 {
			continue
		}
		name := t.Names[0]
		if strings.HasPrefix(name, containerNameNoIndex) {
			matched = append(matched, &t)
		}
	}

	for _, t := range matched {
		containerIndex, err := strconv.Atoi(t.Names[0][len(containerNameNoIndex):])
		if err != nil {
			continue
		}
		if containerIndex >= count {
			slog.Info("Stopping container", slog.String("container_id", t.ID), slog.String("name", t.Names[0]))
			err = c.cli.ContainerStop(context.Background(), t.ID, container.StopOptions{})
			if err != nil {
				slog.Error("Error stopping container", slog.String("container_id", t.ID), slog.String("error", err.Error()))
			}
			err = c.cli.ContainerRemove(context.Background(), t.ID, container.RemoveOptions{})
			if err != nil {
				slog.Error("Error removing container", slog.String("container_id", t.ID), slog.String("error", err.Error()))
			}
		}
	}
}
