package app

import (
	"dockman/app/logger"
	"dockman/app/util/networkutil"
	"fmt"
	"os"
	"runtime"
	"time"
)

func (a *Agent) RegisterMonitor() {
	server, err := ServerGet(a.locator, a.serverId)
	source := "server"
	if err == nil {
		source = fmt.Sprintf("server-%s", server.FormattedName())
	}
	a.registry.GetJobRunner().Add(source, "ServerUpdateStatus", "Sends latest details about the server to the dockman host, the heartbeat.", 3*time.Second, a.updateStatus)
	a.registry.GetJobRunner().Add(source, "ServerResourceStatusMonitor", "Sends latest details about the status of all running resources on the server", 3*time.Second, a.resourceStatusMonitor)
	a.registry.GetJobRunner().Add(source, "ServerMonitorInstanceCount", "Monitors how many resources are currently running vs how many should be based on config and ensures they match.", 3*time.Second, a.monitorInstanceCount)
}

func (a *Agent) resourceStatusMonitor() {
	resources, err := GetResourcesForServer(a.locator, a.serverId)

	if err != nil {
		logger.ErrorWithFields("Failed to get resources for server", err, map[string]any{
			"server_id": a.serverId,
		})
		return
	}

	for _, resource := range resources {
		update, err := a.CalculateResourceServer(resource)

		if err != nil {
			logger.ErrorWithFields("Failed to calculate resource status", err, map[string]any{
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
