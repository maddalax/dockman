package builder

import (
	"github.com/docker/docker/api/types"
	"paas/docker"
	"paas/resources"
)

func (b *ResourceBuilder) runDockerImageBuilder(buildMeta *resources.DockerBuildMeta) error {
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

	err = client.Build(b.OutputStream, result.Directory, types.ImageBuildOptions{
		Dockerfile: buildMeta.Dockerfile,
	})

	if err != nil {
		return b.BuildError(err)
	}

	b.UpdateDeployStatus(resources.DeploymentStatusSucceeded)

	return nil
}
