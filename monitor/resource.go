package monitor

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/docker"
	"paas/domain"
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
				err := resources.SetRunStatus(monitor.locator, res.Id, status)
				if err != nil {
					continue
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func (monitor *Monitor) GetRunStatus(resource *domain.Resource) domain.RunStatus {
	if resource.RunType == domain.RunTypeDockerBuild || resource.RunType == domain.RunTypeDockerRegistry {
		return getRunStatusDocker(resource)
	}
	return domain.RunStatusUnknown
}

func getRunStatusDocker(resource *domain.Resource) domain.RunStatus {
	client, err := docker.Connect()
	if err != nil {
		return domain.RunStatusErrored
	}
	status, err := client.GetRunStatus(resource)
	if err != nil {
		return domain.RunStatusErrored
	}
	return status
}
