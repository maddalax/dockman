package app

import (
	"dockside/app/logger"
	"github.com/maddalax/htmgo/framework/service"
	"time"
)

type LastRunCache[T comparable] struct {
	cache map[string]T
}

func NewLastRunCache[T comparable]() *LastRunCache[T] {
	return &LastRunCache[T]{
		cache: make(map[string]T),
	}
}

// ApplyChange applies a change to the cache, returns true if the value was changed
func (cache *LastRunCache[T]) ApplyChange(id string, value T) bool {
	val, ok := cache.cache[id]
	// value isn't cached, nothing to compare against, return false
	if !ok {
		cache.cache[id] = value
		return false
	}
	if val != value {
		cache.cache[id] = value
		return true
	}
	return false
}

type ResourceMonitor struct {
	locator          *service.Locator
	lastRunStatus    *LastRunCache[RunStatus]
	lastServerStatus *LastRunCache[bool]
}

func NewMonitor(locator *service.Locator) *ResourceMonitor {
	return &ResourceMonitor{
		locator:          locator,
		lastRunStatus:    NewLastRunCache[RunStatus](),
		lastServerStatus: NewLastRunCache[bool](),
	}
}

func (monitor *ResourceMonitor) Start() {
	runner := IntervalJobRunnerFromLocator(monitor.locator)
	runner.Add("ResourceRunStatusMonitor", time.Second*3, monitor.RunStatusMonitorJob)
	runner.Add("ResourceServerCleanup", time.Minute, monitor.ResourceServerCleanup)
	runner.Add("ServerConnectionMonitor", time.Second*5, monitor.ServerConnectionMonitor)
}

// RunStatusMonitorJob Monitors the run status of resources and updates the status if necessary
// Runs every 3s
func (monitor *ResourceMonitor) RunStatusMonitorJob() {
	registry := GetServiceRegistry(monitor.locator)
	list, err := ResourceList(monitor.locator)
	if err != nil {
		return
	}
	for _, res := range list {
		status := GetComputedRunStatus(res)
		changed := monitor.lastRunStatus.ApplyChange(res.Id, status)
		if changed {
			registry.GetEventHandler().OnResourceStatusChange(res, status)
		}
	}
}

// ResourceServerCleanup Cleans up servers that are no longer exist on a resource
func (monitor *ResourceMonitor) ResourceServerCleanup() {
	registry := GetServiceRegistry(monitor.locator)
	list, err := ResourceList(monitor.locator)
	if err != nil {
		logger.Error("Error getting resource list", err)
		return
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
				} else {
					registry.GetEventHandler().OnServerDetached(detail.ServerId, res)
				}
			}
		}
	}
}

// ServerConnectionMonitor Monitors the connection status of servers
func (monitor *ResourceMonitor) ServerConnectionMonitor() {
	registry := GetServiceRegistry(monitor.locator)
	list, err := ServerList(monitor.locator)
	if err != nil {
		logger.Error("Error getting server list", err)
		return
	}
	for _, server := range list {
		accessible := server.IsAccessible()
		changed := monitor.lastServerStatus.ApplyChange(server.Id, accessible)
		if changed {
			if accessible {
				registry.GetEventHandler().OnServerConnected(server)
			} else {
				registry.GetEventHandler().OnServerDisconnected(server)
			}
		}
	}
}
