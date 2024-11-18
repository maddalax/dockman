package builder

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/maddalax/htmgo/framework/h"
	"paas/docker"
	"paas/domain"
	"paas/git"
	"paas/urls"
)

func (b *ResourceBuilder) runDockerImageBuilder(buildMeta *domain.DockerBuildMeta) error {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
	}()

	b.LogBuildMessage("Connecting to Docker...")

	client, err := docker.Connect()

	if err != nil {
		return b.BuildError(err)
	}

	b.UpdateDeployStatus(domain.DeploymentStatusRunning)

	result, err := git.Clone(git.CloneRequest{
		Meta:     buildMeta,
		Progress: b.BuildOutputStream,
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

	err = client.Build(b.BuildOutputStream, result.Directory, types.ImageBuildOptions{
		Dockerfile: buildMeta.Dockerfile,
		BuildID:    dockerBuildId,
		Labels: map[string]string{
			"paas.resource.id": b.Resource.Id,
			"paas.build.id":    b.BuildId,
		},
		Tags: []string{
			fmt.Sprintf(fmt.Sprintf("%s-%s:latest", b.Resource.Name, b.Resource.Id)),
		},
	}, &handlers)

	if err != nil {
		return b.BuildError(err)
	}

	b.LogBuildMessage("Starting container...")

	// build successful, lets try to run it
	err = client.Run(b.Resource, docker.RunOptions{
		RemoveExisting: true,
	})

	if err != nil {
		return b.BuildError(err)
	}

	b.LogBuildMessage("Container successfully started.")
	b.LogBuildMessage(
		h.Render(
			h.A(
				h.Href(urls.ResourceRunLogUrl(b.Resource.Id)),
				h.Text("View run logs"),
				h.Class("underline text-brand-500"),
			),
		),
	)

	b.UpdateDeployStatus(domain.DeploymentStatusSucceeded)

	return nil
}
