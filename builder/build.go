package builder

import (
	"errors"
	"github.com/maddalax/htmgo/framework/service"
	"log/slog"
	"paas/kv"
	"paas/kv/subject"
	"paas/resources"
	"sync"
	"time"
)

var builderLock = sync.Mutex{}

type ResourceBuilder struct {
	Resource           *resources.Resource
	BuildId            string
	ServiceLocator     *service.Locator
	NatsClient         *kv.Client
	OutputStream       *kv.NatsWriter
	LogBuildMessage    func(message string)
	LogBuildError      func(err error)
	UpdateDeployStatus func(status resources.DeploymentStatus)
	PendingCancel      bool
	CancelBuildFunc    func() error
	Finished           bool
}

func NewResourceBuilder(serviceLocator *service.Locator, resource *resources.Resource, buildId string) *ResourceBuilder {
	builderLock.Lock()
	defer builderLock.Unlock()
	// Only want one builder per resource and buildId
	existing := GetBuilder(resource.Id, buildId)
	if existing != nil {
		return existing
	}
	natsClient := service.Get[kv.Client](serviceLocator)
	builder := &ResourceBuilder{
		Resource:       resource,
		BuildId:        buildId,
		NatsClient:     natsClient,
		ServiceLocator: serviceLocator,
		Finished:       false,
		PendingCancel:  false,
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
	SetBuilder(resource.Id, buildId, builder)
	return builder
}

func (b *ResourceBuilder) onFinish() {
	slog.Info("Builder finished", slog.String("resource", b.Resource.Id), slog.String("build", b.BuildId))
	b.Finished = true
	ClearBuilder(b.Resource.Id, b.BuildId)
}

func (b *ResourceBuilder) CancelBuild() {
	b.PendingCancel = true
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

	defer func() {
		b.onFinish()
	}()

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

	go func() {
		for {
			if b.Finished {
				return
			}
			if b.PendingCancel && b.CancelBuildFunc != nil {
				slog.Info("Cancelling build", slog.String("resource", b.Resource.Id), slog.String("build", b.BuildId))
				err := b.CancelBuildFunc()
				if err != nil {
					b.LogBuildError(err)
				}
				b.LogBuildError(errors.New("build cancelled"))
				return
			}
			time.Sleep(time.Millisecond * 250)
		}
	}()

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
