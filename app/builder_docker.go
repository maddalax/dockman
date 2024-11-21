package app

import (
	"context"
	"dockside/app/urls"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/pkg/errors"
	"time"
)

func (b *ResourceBuilder) runDockerImageBuilder(buildMeta *DockerBuildMeta) error {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
	}()

	b.LogBuildMessage("Connecting to Docker...")

	client, err := DockerConnect(b.ServiceLocator)

	if err != nil {
		return b.BuildError(err)
	}

	b.UpdateDeployStatus(DeploymentStatusRunning)

	result, err := buildMeta.CloneRepo(CloneRepoRequest{
		Progress: b.BuildOutputStream,
	})

	if err != nil {
		return b.BuildError(err)
	}

	dockerBuildId := fmt.Sprintf("%s-%s", b.Resource.Id, b.BuildId)

	handlers := BuildResponse{
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

	imageName := fmt.Sprintf("%s-%s", b.Resource.Name, b.Resource.Id)

	err = client.Build(b.BuildOutputStream, result.Directory, types.ImageBuildOptions{
		Dockerfile: buildMeta.Dockerfile,
		BuildID:    dockerBuildId,
		Labels: map[string]string{
			"dockside.resource.id": b.Resource.Id,
			"dockside.build.id":    b.BuildId,
		},
		Tags: []string{
			fmt.Sprintf(fmt.Sprintf("%s:latest", imageName)),
			fmt.Sprintf(fmt.Sprintf("%s:buildId-%s", imageName, b.BuildId)),
		},
	}, &handlers)

	if err != nil {
		return b.BuildError(err)
	}

	b.LogBuildMessage("Saving image...")

	err = client.SaveImage(imageName, b.BuildId)

	if err != nil {
		return b.BuildError(err)
	}

	b.UpdateDeployStatus(DeploymentStatusSucceeded)

	b.LogBuildMessage("Successfully saved image, starting process on enabled servers...")

	responses, err := SendCommandForResource[RunResourceResponse](b.ServiceLocator, b.Resource.Id, SendCommandOpts{
		Command: &RunResourceCommand{
			ResourceId: b.Resource.Id,
		},
		Timeout: time.Second * 5,
	})

	if err != nil {
		return b.BuildError(err)
	}

	for _, response := range responses {
		if response.Response.Error != nil {
			b.LogBuildError(errors.Wrap(response.Response.Error,
				fmt.Sprintf("Failed to start resource on server %s", response.ServerDetails.Hostname)))
		} else {
			b.LogBuildMessage(fmt.Sprintf("Resource started on server %s", response.ServerDetails.Hostname))
		}
	}

	b.LogBuildMessage(
		h.Render(
			h.A(
				h.Href(urls.ResourceRunLogUrl(b.Resource.Id)),
				h.Text("View run logs"),
				h.Class("underline text-brand-500"),
			),
		),
	)

	return nil
}
