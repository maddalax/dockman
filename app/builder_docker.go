package app

import (
	"context"
	"dockside/app/urls"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/pkg/errors"
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
		UseCache:     false,
		Progress:     b.BuildOutputStream,
		SingleBranch: true,
		BranchName:   buildMeta.DeploymentBranch,
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
			"git.commit.hash":      result.Commit,
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

	b.LogBuildMessage(fmt.Sprintf("Container built with commit %s", result.Commit))

	b.UpdateDeployStatus(DeploymentStatusSucceeded)

	err = ResourcePatch(b.ServiceLocator, b.Resource.Id, func(resource *Resource) *Resource {
		resource.BuildMeta.(*DockerBuildMeta).CommitForBuild = result.Commit
		return resource
	})

	if err != nil {
		b.LogBuildError(err)
	}

	b.LogBuildMessage("Successfully saved image, starting process on enabled servers...")

	responses, err := SendResourceStartCommand(b.ServiceLocator, b.Resource.Id, StartOpts{
		RemoveExisting: true,
	})

	if err != nil {
		return b.BuildError(err)
	}

	for _, response := range responses {

		serverName := h.Ternary(response.ServerDetails.Name == "", response.ServerDetails.HostName, response.ServerDetails.Name)

		if response.Response.Error != "" || response.SendError != nil {
			err = h.Ternary(response.Response.Error == "", response.SendError, errors.New(response.Response.Error)).(error)
			b.LogBuildError(errors.Wrap(err, fmt.Sprintf("Failed to start resource on server %s", serverName)))
		} else {
			b.LogBuildMessage(fmt.Sprintf("Resource started on server %s", serverName))
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
