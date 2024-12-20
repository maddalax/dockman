package app

import (
	"dockman/app/logger"
	"dockman/app/subject"
	"github.com/maddalax/htmgo/framework/service"
	"sync"
	"time"
)

var builderLock = sync.Mutex{}

type ResourceBuilder struct {
	Source            string
	Resource          *Resource
	BuilderRegistry   *BuilderRegistry
	BuildId           string
	ServiceLocator    *service.Locator
	NatsClient        *KvClient
	BuildOutputStream *NatsWriter
	// RunOutputStream stream for run logs that we don't want to persist since docker already persists them
	RunOutputStream    *EphemeralNatsWriter
	LogBuildMessage    func(message string)
	LogBuildError      func(err error)
	LogRunMessage      func(message string)
	UpdateDeployStatus func(status DeploymentStatus)
	PatchDeployment    func(cb func(deployment *Deployment) *Deployment)
	PendingCancel      bool
	CancelBuildFunc    func() error
	Finished           bool
}

func NewResourceBuilder(serviceLocator *service.Locator, resource *Resource, buildId string, source string) *ResourceBuilder {
	builderLock.Lock()
	defer builderLock.Unlock()
	registry := GetBuilderRegistry(serviceLocator)
	// Only want one builder per resource and buildId
	existing := registry.GetBuilder(resource.Id, buildId)
	if existing != nil {
		return existing
	}
	natsClient := service.Get[KvClient](serviceLocator)
	builder := &ResourceBuilder{
		Source:          source,
		BuilderRegistry: registry,
		Resource:        resource,
		BuildId:         buildId,
		NatsClient:      natsClient,
		ServiceLocator:  serviceLocator,
		Finished:        false,
		PendingCancel:   false,
		LogBuildMessage: func(message string) {
			natsClient.LogBuildMessage(resource.Id, buildId, message)
		},
		LogBuildError: func(err error) {
			natsClient.LogBuildError(resource.Id, buildId, err)
		},
		LogRunMessage: func(message string) {
			natsClient.LogRunMessage(resource.Id, message)
		},
		UpdateDeployStatus: func(status DeploymentStatus) {
			err := PatchDeployment(serviceLocator, resource.Id, buildId, func(deployment *Deployment) *Deployment {
				deployment.Status = status
				return deployment
			})
			if err != nil {
				logger.ErrorWithFields("failed to update deployment status", err, map[string]any{
					"resource": resource.Id,
					"build":    buildId,
					"status":   status,
				})
			}
		},
		PatchDeployment: func(cb func(deployment *Deployment) *Deployment) {
			err := PatchDeployment(serviceLocator, resource.Id, buildId, cb)
			if err != nil {
				logger.ErrorWithFields("failed to patch deployment", err, map[string]any{
					"resource": resource.Id,
					"build":    buildId,
				})
			}
		},
	}
	registry.SetBuilder(resource.Id, buildId, builder)
	return builder
}

func (b *ResourceBuilder) CanBuild() error {
	if len(b.Resource.ServerDetails) == 0 {
		return NoServersAttachedError
	}
	return nil
}

func (b *ResourceBuilder) onFinish() {
	logger.InfoWithFields("Builder finished", map[string]any{
		"resource": b.Resource.Id,
		"build":    b.BuildId,
	})
	b.Finished = true
	b.BuilderRegistry.ClearBuilder(b.Resource.Id, b.BuildId)
}

func (b *ResourceBuilder) CancelBuild() {
	b.PendingCancel = true
}

func (b *ResourceBuilder) ClearLogs() {
	err := b.NatsClient.PurgeStream(
		b.NatsClient.BuildLogStreamName(b.Resource.Id, b.BuildId),
	)
	if err != nil {
		logger.ErrorWithFields("failed to clear build logs", err, map[string]any{
			"resource": b.Resource.Id,
			"build":    b.BuildId,
		})
	}
}

func (b *ResourceBuilder) BuildError(err error) error {
	b.LogBuildError(err)
	b.PatchDeployment(func(deployment *Deployment) *Deployment {
		deployment.Status = DeploymentStatusFailed
		deployment.StatusReason = err.Error()
		return deployment
	})
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

	err = b.CanBuild()

	if err != nil {
		return b.BuildError(err)
	}

	err = b.NatsClient.CreateRunLogStream(b.Resource.Id)

	if err != nil {
		return err
	}

	b.BuildOutputStream = b.NatsClient.NewNatsWriter(subject.BuildLogForResource(b.Resource.Id, b.BuildId))
	b.RunOutputStream = b.NatsClient.NewEphemeralNatsWriter(subject.RunLogsForResource(b.Resource.Id))

	err = CreateDeployment(b.ServiceLocator, CreateDeploymentRequest{
		ResourceId: b.Resource.Id,
		BuildId:    b.BuildId,
		Source:     b.Source,
	})

	if err != nil {
		return b.BuildError(err)
	}

	b.UpdateDeployStatus(DeploymentStatusPending)

	go func() {
		for {
			if b.Finished {
				return
			}
			if b.PendingCancel && b.CancelBuildFunc != nil {
				logger.InfoWithFields("Cancelling build", map[string]any{
					"resource": b.Resource.Id,
					"build":    b.BuildId,
				})
				err := b.CancelBuildFunc()
				if err != nil {
					b.LogBuildError(err)
				}
				b.LogBuildError(BuildCancelledError)
				return
			}
			time.Sleep(time.Millisecond * 250)
		}
	}()

	switch bm := b.Resource.BuildMeta.(type) {
	case *DockerBuildMeta:
		return b.runDockerImageBuilder(bm)
	default:
		return UnknownBuildTypeError
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
