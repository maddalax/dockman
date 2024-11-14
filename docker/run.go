package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/go-connections/nat"
	"io"
	"paas/resources"
	"strconv"
)

type RunOptions struct {
	Stdout io.Writer
	// If we should kill the existing container that's running first
	KillExisting bool
}

func (c *Client) GetContainer(resource *resources.Resource) (types.ContainerJSON, error) {
	containerName := fmt.Sprintf("%s-%s-container", resource.Name, resource.Id)
	return c.cli.ContainerInspect(context.Background(), containerName)
}

func (c *Client) Run(resource *resources.Resource, opts RunOptions) error {
	ctx := context.Background()
	imageName := fmt.Sprintf("%s-%s", resource.Name, resource.Id)
	containerName := fmt.Sprintf("%s-%s-container", resource.Name, resource.Id)

	err := c.cli.ContainerStop(ctx, containerName, container.StopOptions{})

	if err != nil {
		switch err.(type) {
		case errdefs.ErrNotFound:
			// don't need to worry about it if the container doesn't exist
			err = nil
		default:
			return err
		}
	}

	err = c.cli.ContainerRemove(ctx, containerName, container.RemoveOptions{
		Force: true,
	})

	if err != nil {
		switch err.(type) {
		case errdefs.ErrNotFound:
			// don't need to worry about it if the container doesn't exist
			err = nil
		default:
			return err
		}
	}

	hostPort, err := FindOpenPort(3000)

	if err != nil {
		return err
	}

	// Define port bindings
	portBindings := nat.PortMap{
		"3000/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",              // Bind to all network interfaces
				HostPort: strconv.Itoa(hostPort), // Map container port 80 to host port 8080
			},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		LogConfig: container.LogConfig{
			Type: "json-file",
			Config: map[string]string{
				"max-size": "10m",
			},
		},
	}

	// Create and start a container
	resp, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		ExposedPorts: map[nat.Port]struct{}{
			// the port the container exposes
			"3000/tcp": {},
		},
		AttachStdout: true,
		AttachStderr: true,
	}, hostConfig, nil, nil, containerName)

	if err != nil {
		return err
	}

	if err := c.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	if opts.Stdout != nil {
		return c.StreamLogs(resp.ID, StreamLogsOptions{
			Stdout: opts.Stdout,
		})
	}

	return nil
}
