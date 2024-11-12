package builder

import (
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/maddalax/htmgo/framework/h"
	"paas/docker"
	"paas/kv"
	"paas/kv/subject"
	"paas/resources"
	"time"
)

func StartBuildAsync(ctx *h.RequestContext, resource *resources.Resource, buildId string, wait time.Duration) error {
	// just ensure we create the stream, and then start the build in the background
	natsClient := kv.GetClientFromCtx(ctx)

	err := natsClient.CreateBuildLogStream(resource.Id, buildId)

	if err != nil {
		return err
	}

	natsClient.LogBuildMessage(resource.Id, buildId, "Starting build...")

	go func() {
		time.Sleep(wait)
		err := BuildResource(ctx, resource, buildId)
		if err != nil {
			natsClient.LogBuildError(resource.Id, buildId, err)
		}
	}()

	return nil
}

func BuildResource(ctx *h.RequestContext, resource *resources.Resource, buildId string) error {
	natsClient := kv.GetClientFromCtx(ctx)

	err := natsClient.CreateBuildLogStream(resource.Id, buildId)

	if err != nil {
		return err
	}

	outputStream := natsClient.NewNatsWriter(subject.BuildLogForResource(resource.Id, buildId))

	err = resources.CreateDeployment(ctx.ServiceLocator(), resources.CreateDeploymentRequest{
		ResourceId: resource.Id,
		BuildId:    buildId,
	})

	if err != nil {
		return err
	}

	outputStream.Write([]byte("Sydne is a fart\n"))

	natsClient.LogBuildMessage(resource.Id, buildId, "Connecting to Docker...")

	client, err := docker.Connect()

	if err != nil {
		return err
	}

	switch bm := resource.BuildMeta.(type) {
	case *resources.DockerBuildMeta:
		return buildDocker(client, bm, outputStream)
	default:
		return errors.New("unknown build type")
	}
}

func buildDocker(client *docker.Client, buildMeta *resources.DockerBuildMeta, outputStream *kv.NatsWriter) error {
	result, err := resources.Clone(resources.CloneRequest{
		Meta:     buildMeta,
		Progress: outputStream,
	})
	if err != nil {
		return err
	}
	err = client.Build(outputStream, result.Directory, types.ImageBuildOptions{
		Dockerfile: buildMeta.Dockerfile,
	})
	if err != nil {
		return err
	}
	return nil
}
