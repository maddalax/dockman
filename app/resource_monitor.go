package app

import (
	"dockman/app/logger"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"slices"
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
	source := "dockman"
	runner.Add(source, "ResourceRunStatusMonitor", "Checks to see if resources are running or stopped", time.Second*3, monitor.RunStatusMonitorJob)
	runner.Add(source, "ResourceServerCleanup", "Detaches servers that no longer exist from resources", time.Minute, monitor.ResourceServerCleanup)
	runner.Add(source, "ServerConnectionMonitor", "Monitors if connected servers are still connected by checking for a heartbeat", time.Second*5, monitor.ServerConnectionMonitor)
	runner.Add(source, "ResourceCheckForNewCommits", "Checks if a resource has a new commit and starts a new deployment if enabled", time.Second*30, monitor.ResourceCheckForNewCommits)
	runner.Add(source, "ServerDuplicateCleanup", "Checks if there are any servers with the same remote ip and dedupes them", time.Second*30, monitor.CleanupDuplicateServers)

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

func (monitor *ResourceMonitor) ResourceCheckForNewCommits() {
	registry := GetServiceRegistry(monitor.locator)
	list, err := ResourceList(monitor.locator)
	if err != nil {
		logger.Error("Error getting resource list", err)
		return
	}
	for _, res := range list {
		switch bm := res.BuildMeta.(type) {
		case *DockerBuildMeta:
			if !bm.DeployOnNewCommit {
				continue
			}
			latest, err := bm.GetLatestCommitOnRemote()
			if err != nil {
				logger.ErrorWithFields("Error getting latest commit", err, map[string]interface{}{
					"resource": res.Id,
				})
				continue
			}
			current := bm.CommitForBuild
			logger.DebugWithFields("Checking for new commits", map[string]interface{}{
				"resource": res.Id,
				"latest":   latest,
				"current":  current,
			})
			if current != "" && latest != "" && latest != current {
				registry.GetEventHandler().OnNewCommit(res, bm.DeploymentBranch, latest)
			}
		}
	}
}

// CleanupDuplicateServers Detaches and deletes servers that have the same remote ip, keeping the newest one
// this can happen if a server os is reinstalled and has a new id
func (monitor *ResourceMonitor) CleanupDuplicateServers() {
	registry := GetServiceRegistry(monitor.locator)
	list, err := ServerList(monitor.locator)
	if err != nil {
		logger.Error("Error getting server list", err)
		return
	}
	serverMap := make(map[string][]*Server)
	for _, server := range list {
		ip := server.IpAddress()
		if serverMap[ip] == nil {
			serverMap[ip] = []*Server{}
		}
		serverMap[ip] = append(serverMap[ip], server)
	}

	for _, servers := range serverMap {
		if len(servers) < 2 {
			continue
		}

		// sort by last seen, newest first
		slices.SortFunc(servers, func(a, b *Server) int {
			return b.LastSeen.Compare(a.LastSeen)
		})

		// check the rest of the servers, skip the first one because it's the newest
		inaccessibleServers := h.Filter(servers[1:], func(s *Server) bool {
			return !s.IsAccessible()
		})

		// all servers accessible, this really shouldn't happen... but safety first
		if len(inaccessibleServers) == 0 {
			logger.WarnWithFields("All servers are accessible, not detaching any", map[string]interface{}{
				"server_ids": h.Map(servers, func(s *Server) string { return s.Id }),
			})
			continue
		}

		// all servers are inaccessible, don't detach any
		if len(inaccessibleServers) == len(servers) {
			logger.DebugWithFields("All servers are inaccessible, not detaching any", map[string]interface{}{
				"server_ids": h.Map(servers, func(s *Server) string { return s.Id }),
			})
			continue
		}

		// detach from all resources
		for _, server := range inaccessibleServers {
			logger.InfoWithFields("Detaching duplicate server", map[string]interface{}{
				"server_id": server.Id,
			})
			resources, err := GetResourcesForServer(monitor.locator, server.Id)
			if err != nil {
				continue
			}
			for _, resource := range resources {
				err := DetachServerFromResource(monitor.locator, server.Id, resource.Id)
				if err != nil {
					logger.Error("Error detaching duplicate server", err)
				} else {
					registry.GetEventHandler().OnServerDetached(server.Id, resource)
				}
			}
		}

		for _, server := range inaccessibleServers {
			logger.InfoWithFields("Deleting duplicate server", map[string]interface{}{
				"server_id": server.Id,
			})
			err := ServerDelete(monitor.locator, server.Id)
			if err != nil {
				logger.Error("Error deleting duplicate server", err)
			}
		}
	}

}
