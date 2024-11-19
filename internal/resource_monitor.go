package internal

import (
	"context"
	"github.com/maddalax/htmgo/framework/service"
	"paas/internal/subject"
	"time"
)

type ResourceMonitor struct {
	locator *service.Locator
}

func NewMonitor(locator *service.Locator) *ResourceMonitor {
	return &ResourceMonitor{
		locator: locator,
	}
}

func (monitor *ResourceMonitor) StartRunStatusMonitor() {
	for {
		list, err := ResourceList(monitor.locator)
		if err != nil {
			continue
		}
		for _, res := range list {
			status := monitor.GetRunStatus(res)
			if res.RunStatus != status {
				err := SetRunStatus(monitor.locator, res.Id, status)
				if err != nil {
					continue
				}
				monitor.OnStatusChange(res, status)
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func (monitor *ResourceMonitor) OnStatusChange(resource *Resource, status RunStatus) {
	ctx, cancel := context.WithCancel(context.Background())
	natsClient := KvFromLocator(monitor.locator)
	writer := natsClient.CreateEphemeralWriterSubscriber(ctx, subject.RunLogsForResource(resource.Id), NatsWriterCreateOptions{})

	message := ""
	if status == RunStatusRunning {
		message = "Container is now running"
	} else {
		message = "Container has stopped"
	}

	writer.Writer.Write([]byte(message))

	cancel()
}

func (monitor *ResourceMonitor) GetRunStatus(resource *Resource) RunStatus {
	if resource.RunType == RunTypeDockerBuild || resource.RunType == RunTypeDockerRegistry {
		return getRunStatusDocker(resource)
	}
	return RunStatusUnknown
}

func getRunStatusDocker(resource *Resource) RunStatus {
	client, err := DockerConnect()
	if err != nil {
		return RunStatusNotRunning
	}
	status, err := client.GetRunStatus(resource)
	if err != nil {
		return RunStatusNotRunning
	}
	return status
}
