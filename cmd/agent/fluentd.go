package main

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
	"paas/app"
	"path/filepath"
	"strings"
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

	fluentConfPath, err := filepath.Abs("./fluent.conf")

	if err != nil {
		return err
	}

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

	// Create the container
	resp, err := cli.ContainerCreate(context.Background(), config, hostConfig, networkConfig, nil, m.containerName)

	if err != nil {
		// container already running, just return
		if strings.Contains(err.Error(), "already in use") {
			return nil
		}
		return err
	}

	// Start the container
	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	return nil
}

func (m *FluentdManager) StreamLogs() error {
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

	resp, err := cli.ContainerExecAttach(ctx, execIDResp.ID, container.ExecAttachOptions{})

	if err != nil {
		return err
	}

	defer resp.Close()

	// Read and print the output
	reader := bufio.NewReader(resp.Reader)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		m.agent.WriteContainerLog(line)
	}

	return nil
}
