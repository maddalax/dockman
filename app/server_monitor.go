package app

import (
	"dockside/app/logger"
	"dockside/app/util/networkutil"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"os"
	"runtime"
	"strings"
	"time"
)

func (a *Agent) RegisterMonitor() {
	server, err := ServerGet(a.locator, a.serverId)
	name := "unknown"
	if err == nil {
		name = server.FormattedName()
	}
	a.registry.GetJobRunner().Add(fmt.Sprintf("ServerUpdateStatus-%s-%s", a.serverId, name), 3*time.Second, a.updateStatus)
	a.registry.GetJobRunner().Add(fmt.Sprintf("ResourceStatusMonitor-%s-%s", a.serverId, name), 3*time.Second, a.resourceStatusMonitor)
}

func (a *Agent) resourceStatusMonitor() {
	resources, err := GetResourcesForServer(a.locator, a.serverId)

	if err != nil {
		logger.ErrorWithFields("Failed to get resources for server", err, map[string]any{
			"server_id": a.serverId,
		})
		return
	}

	logger.InfoWithFields("Updating resource statuses", map[string]any{
		"server_id": a.serverId,
		"count":     len(resources),
		"resource_ids": strings.Join(h.Map(resources, func(r *Resource) string {
			return r.Id
		}), ", "),
	})

	for _, resource := range resources {
		update, err := a.CalculateResourceServer(resource)

		if err != nil {
			logger.ErrorWithFields("Failed to calculate resource status", nil, map[string]any{
				"resource_id": resource.Id,
			})
			continue
		}

		err = PatchResourceServer(a.locator, resource.Id, a.serverId, func(server *ResourceServer) *ResourceServer {
			server.Upstreams = update.Upstreams
			server.RunStatus = update.RunStatus
			server.LastUpdate = time.Now()
			return server
		})
		logger.InfoWithFields("Updated resource status", map[string]any{
			"resource_id": resource.Id,
			"server_id":   a.serverId,
			"status":      update.RunStatus,
			"upstreams":   len(update.Upstreams),
		})
		if err != nil {
			logger.ErrorWithFields("Failed to update resource status", err, map[string]any{
				"resource_id": resource.Id,
				"server_id":   a.serverId,
				"status":      update.RunStatus,
				"upstreams":   len(update.Upstreams),
			})
		}
	}
}

func (a *Agent) updateStatus() {
	hostName, err := os.Hostname()

	if err != nil {
		hostName = ""
	}

	localIp := networkutil.GetLocalIp()

	err = ServerPut(a.locator, ServerPutOpts{
		Id:              a.serverId,
		HostName:        hostName,
		LocalIpAddress:  localIp,
		RemoteIpAddress: "",
		LastSeen:        time.Now(),
		Os:              fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH),
	})

	if err != nil {
		logger.ErrorWithFields("Failed to update server status", err, map[string]any{
			"server_id": a.serverId,
		})
	}
}
