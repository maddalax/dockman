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

func (a *Agent) StartServerMonitor() {
	for {
		a.updateStatus()
		time.Sleep(3 * time.Second)
	}
}

func (a *Agent) StartResourceStatusMonitor() {
	for {
		time.Sleep(3 * time.Second)
		resources, err := GetResourcesForServer(a.locator, a.serverId)
		if err != nil {
			logger.ErrorWithFields("Failed to get resources for server", err, map[string]any{
				"server_id": a.serverId,
			})
			continue
		}

		logger.InfoWithFields("Updating resource statuses", map[string]any{
			"server_id": a.serverId,
			"count":     len(resources),
			"resource_ids": strings.Join(h.Map(resources, func(r *Resource) string {
				return r.Id
			}), ", "),
		})

		for _, resource := range resources {
			status := a.GetRunStatus(resource)
			err := PatchResourceServer(a.locator, resource.Id, a.serverId, func(server *ResourceServer) *ResourceServer {
				server.RunStatus = status
				return server
			})
			logger.InfoWithFields("Updated resource status", map[string]any{
				"resource_id": resource.Id,
				"server_id":   a.serverId,
				"status":      status,
			})
			if err != nil {
				logger.ErrorWithFields("Failed to update resource status", err, map[string]any{
					"resource_id": resource.Id,
					"server_id":   a.serverId,
					"status":      status,
				})
			}
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
