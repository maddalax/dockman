package builder

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"paas/docker"
	"paas/resources"
)

func (b *ResourceBuilder) runDockerImageBuilder(buildMeta *resources.DockerBuildMeta) error {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
	}()

	b.LogBuildMessage("Connecting to Docker...")

	client, err := docker.Connect()

	if err != nil {
		return b.BuildError(err)
	}

	b.UpdateDeployStatus(resources.DeploymentStatusRunning)

	result, err := resources.Clone(resources.CloneRequest{
		Meta:     buildMeta,
		Progress: b.OutputStream,
	})

	if err != nil {
		return b.BuildError(err)
	}

	dockerBuildId := fmt.Sprintf("%s-%s", b.Resource.Id, b.BuildId)

	handlers := docker.BuildResponse{
		CancelChan: make(chan func() error),
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		case f := <-handlers.CancelChan:
			b.CancelBuildFunc = f
			return
		}
	}()

	err = client.Build(b.OutputStream, result.Directory, types.ImageBuildOptions{
		Dockerfile: buildMeta.Dockerfile,
		BuildID:    dockerBuildId,
		Labels: map[string]string{
			"paas.resource.id": b.Resource.Id,
			"paas.build.id":    b.BuildId,
		},
		Tags: []string{
			fmt.Sprintf(fmt.Sprintf("%s:latest", b.Resource.Name)),
		},
	}, &handlers)

	if err != nil {
		return b.BuildError(err)
	}

	b.UpdateDeployStatus(resources.DeploymentStatusSucceeded)

	return nil
}
