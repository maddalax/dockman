package main

import (
	"context"
	"dockside/app"
	"dockside/app/logger"
	"dockside/app/volume"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var FluentdConfig = `
<source>
  @type forward
  port 24224
</source>


<match **>
  @type file
  path /fluentd/log/docker.log
  append true
  <format>
    @type json
  </format>
  <buffer>
    @type file
    path /fluentd/log/buffer
    flush_mode interval
    flush_interval 1s
    flush_thread_count 1
    retry_type exponential_backoff
    retry_wait 1s
    retry_max_interval 30s
    retry_forever true
  </buffer>
</match>
`

type FluentdManager struct {
	agent         *app.Agent
	containerName string
}

func NewFluentdManager(agent *app.Agent) *FluentdManager {
	return &FluentdManager{
		agent:         agent,
		containerName: "fluentd-agent",
	}
}

func (m *FluentdManager) StartContainer() error {
	// Start the fluentd agent
	env := client.FromEnv
	cli, err := client.NewClientWithOpts(env,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return err
	}

	imageName := "fluent/fluentd:v1.17-debian-1"

	out, err := cli.ImagePull(context.Background(), imageName, image.PullOptions{})

	if err != nil {
		return err
	}

	// Display the pull output
	_, err = io.Copy(io.Discard, out) // Use io.Discard to avoid verbose output

	if err != nil {
		return err
	}

	// Container configuration
	config := &container.Config{
		Image: imageName,
		ExposedPorts: map[nat.Port]struct{}{
			"24224/tcp": {},
			"24224/udp": {},
		},
	}

	fluentConfPath := filepath.Join(volume.GetPersistentVolumePath(), "fluentd.conf")

	err = os.WriteFile(fluentConfPath, []byte(FluentdConfig), 0644)

	if err != nil {
		return err
	}

	// Host configuration
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"24224/tcp": []nat.PortBinding{{HostPort: "24224"}},
			"24224/udp": []nat.PortBinding{{HostPort: "24224"}},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: fluentConfPath,
				Target: "/fluentd/etc/fluent.conf",
			},
		},
	}

	// Network configuration
	networkConfig := &network.NetworkingConfig{}

	// Remove the container if it already exists
	err = cli.ContainerRemove(context.Background(), m.containerName, container.RemoveOptions{})

	if err != nil {
		// container is already running successfully so just return
		if strings.Contains(err.Error(), "container is running") {
			return nil
		}
		if !client.IsErrNotFound(err) || strings.Contains(err.Error(), "No such container") {
			err = nil
		} else {
			return err
		}
	}

	// Create the container
	_, err = cli.ContainerCreate(context.Background(), config, hostConfig, networkConfig, nil, m.containerName)

	if err != nil {
		return err
	}

	// Start the container
	if err := cli.ContainerStart(context.Background(), m.containerName, container.StartOptions{}); err != nil {
		return err
	}

	return nil
}

func (m *FluentdManager) StreamLogs() error {
	for {
		env := client.FromEnv
		cli, err := client.NewClientWithOpts(env,
			client.WithAPIVersionNegotiation(),
		)
		if err != nil {
			return err
		}

		execConfig := container.ExecOptions{
			Cmd:          strslice.StrSlice{"bin/bash", "-c", "cd /fluentd/log && find . -type f -name \"docker.log*\" -exec tail -n 0 -f {} +\n"},
			AttachStdout: true,
			AttachStderr: true,
		}

		ctx := context.Background()

		execIDResp, err := cli.ContainerExecCreate(ctx, m.containerName, execConfig)

		if err != nil {
			return err
		}

		logger.InfoWithFields("attaching to fluentd to stream logs", map[string]any{
			"container": m.containerName,
			"exec_id":   execIDResp.ID,
		})
		resp, err := cli.ContainerExecAttach(ctx, execIDResp.ID, container.ExecAttachOptions{})

		if err != nil {
			return err
		}

		server, _ := app.ServerGet(m.agent.GetLocator(), m.agent.GetServerId())

		writer := app.NewNatsContainerLogWriter(m.agent.GetLocator(), server, m.agent)

		_, err = stdcopy.StdCopy(writer, writer, resp.Reader)

		if err != nil {
			logger.ErrorWithFields("failed to stream logs", err, map[string]any{
				"container": m.containerName,
			})
		} else {
			logger.InfoWithFields("finished streaming logs", map[string]any{
				"container": m.containerName,
			})
		}

		// If we disconnect, re-stream logs
		time.Sleep(3 * time.Second)
	}
}
