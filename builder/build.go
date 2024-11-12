package builder

import (
	"errors"
	"github.com/maddalax/htmgo/framework/service"
	"log/slog"
	"paas/kv"
	"paas/kv/subject"
	"paas/resources"
	"time"
)

type ResourceBuilder struct {
	Resource           *resources.Resource
	BuildId            string
	ServiceLocator     *service.Locator
	NatsClient         *kv.Client
	OutputStream       *kv.NatsWriter
	LogBuildMessage    func(message string)
	LogBuildError      func(err error)
	UpdateDeployStatus func(status resources.DeploymentStatus)
}

func NewResourceBuilder(serviceLocator *service.Locator, resource *resources.Resource, buildId string) *ResourceBuilder {
	natsClient := service.Get[kv.Client](serviceLocator)
	return &ResourceBuilder{
		Resource:       resource,
		BuildId:        buildId,
		NatsClient:     natsClient,
		ServiceLocator: serviceLocator,
		LogBuildMessage: func(message string) {
			natsClient.LogBuildMessage(resource.Id, buildId, message)
		},
		LogBuildError: func(err error) {
			natsClient.LogBuildError(resource.Id, buildId, err)
		},
		UpdateDeployStatus: func(status resources.DeploymentStatus) {
			err := resources.UpdateDeploymentStatus(serviceLocator, resources.UpdateDeploymentStatusRequest{
				ResourceId: resource.Id,
				BuildId:    buildId,
				Status:     status,
			})
			if err != nil {
				slog.Error("failed to update deployment status",
					slog.String("resource", resource.Id),
					slog.String("build", buildId),
					slog.String("status", string(status)),
					slog.String("error", err.Error()))
			}
		},
	}
}

func (b *ResourceBuilder) ClearLogs() {
	err := b.NatsClient.PurgeStream(
		b.NatsClient.BuildLogStreamName(b.Resource.Id, b.BuildId),
	)
	if err != nil {
		slog.Error("failed to clear build logs",
			slog.String("resource", b.Resource.Id),
			slog.String("build", b.BuildId),
			slog.String("error", err.Error()))
	}
}

func (b *ResourceBuilder) BuildError(err error) error {
	b.LogBuildError(err)
	b.UpdateDeployStatus(resources.DeploymentStatusFailed)
	return err
}

func (b *ResourceBuilder) Build() error {
	err := b.NatsClient.CreateBuildLogStream(b.Resource.Id, b.BuildId)

	if err != nil {
		return err
	}

	b.OutputStream = b.NatsClient.NewNatsWriter(subject.BuildLogForResource(b.Resource.Id, b.BuildId))

	err = resources.CreateDeployment(b.ServiceLocator, resources.CreateDeploymentRequest{
		ResourceId: b.Resource.Id,
		BuildId:    b.BuildId,
	})

	if err != nil {
		return b.BuildError(err)
	}

	b.UpdateDeployStatus(resources.DeploymentStatusPending)

	switch bm := b.Resource.BuildMeta.(type) {
	case *resources.DockerBuildMeta:
		return b.runDockerImageBuilder(bm)
	default:
		return errors.New("unknown build type")
	}
}

func (b *ResourceBuilder) StartBuildAsync(wait time.Duration) error {
	err := b.NatsClient.CreateBuildLogStream(b.Resource.Id, b.BuildId)

	if err != nil {
		return err
	}

	b.LogBuildMessage("Starting build...")

	go func() {
		time.Sleep(wait)
		// b.Build does its own error handling
		_ = b.Build()
	}()

	return nil
}
