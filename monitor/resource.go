package monitor

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/docker"
	"paas/resources"
	"time"
)

type Monitor struct {
	locator *service.Locator
}

func NewMonitor(locator *service.Locator) *Monitor {
	return &Monitor{
		locator: locator,
	}
}

func (monitor *Monitor) StartRunStatusMonitor() {
	for {
		list, err := resources.List(monitor.locator)
		if err != nil {
			continue
		}
		for _, res := range list {
			status := monitor.GetRunStatus(res)
			if res.RunStatus != status {
				err := res.Patch(monitor.locator, map[string]any{
					"run_status": status,
				})
				if err != nil {
					continue
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func (monitor *Monitor) GetRunStatus(resource *resources.Resource) resources.RunStatus {
	if resource.RunType == resources.RunTypeDockerBuild || resource.RunType == resources.RunTypeDockerRegistry {
		return getRunStatusDocker(resource)
	}
	return resources.RunStatusUnknown
}

func getRunStatusDocker(resource *resources.Resource) resources.RunStatus {
	client, err := docker.Connect()
	if err != nil {
		return resources.RunStatusErrored
	}
	status, err := client.GetRunStatus(resource)
	if err != nil {
		return resources.RunStatusErrored
	}
	return status
}
