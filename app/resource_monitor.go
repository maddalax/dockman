package app

import (
	"context"
	"dockside/app/logger"
	"dockside/app/subject"
	"github.com/maddalax/htmgo/framework/service"
	"time"
)

type ResourceMonitor struct {
	locator       *service.Locator
	lastRunStatus map[string]RunStatus
}

func NewMonitor(locator *service.Locator) *ResourceMonitor {
	return &ResourceMonitor{
		locator:       locator,
		lastRunStatus: make(map[string]RunStatus),
	}
}

func (monitor *ResourceMonitor) Start() {
	go monitor.StartRunStatusMonitor()
	go monitor.StartResourceServerCleanup()
}

// StartRunStatusMonitor Monitors the run status of resources and updates the status if necessary
// Runs every 3s
func (monitor *ResourceMonitor) StartRunStatusMonitor() {
	for {
		list, err := ResourceList(monitor.locator)
		if err != nil {
			continue
		}
		for _, res := range list {
			status := GetComputedRunStatus(res)
			lastStatus, ok := monitor.lastRunStatus[res.Id]
			if !ok {
				monitor.lastRunStatus[res.Id] = status
				continue
			}
			if lastStatus != status {
				monitor.OnStatusChange(res, status)
			}
			monitor.lastRunStatus[res.Id] = status
		}
		time.Sleep(3 * time.Second)
	}
}

// StartResourceServerCleanup Cleans up servers that are no longer exist on a resource
// Runs every 60s
// TODO have some way to monitor these jobs
func (monitor *ResourceMonitor) StartResourceServerCleanup() {
	for {
		list, err := ResourceList(monitor.locator)
		if err != nil {
			time.Sleep(time.Second)
			logger.Error("Error getting resource list", err)
			continue
		}
		for _, res := range list {
			for _, detail := range res.ServerDetails {
				_, err := ServerGet(monitor.locator, detail.ServerId)
				if err != nil && err.Error() == NatsKeyNotFoundError.Error() {
					logger.WarnWithFields("server no longer exists, detaching it from resource", map[string]interface{}{
						"server_id":   detail.ServerId,
						"resource_id": res.Id,
					})
					err := DetachServerFromResource(monitor.locator, detail.ServerId, res.Id)
					if err != nil {
						logger.ErrorWithFields("Error detaching server from resource", err, map[string]interface{}{
							"server_id":   detail.ServerId,
							"resource_id": res.Id,
						})
					}
				}
			}
		}
		time.Sleep(time.Minute)
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
