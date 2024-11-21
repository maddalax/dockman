package app

import (
	"dockside/app/logger"
	"dockside/app/util/networkutil"
	"fmt"
	"os"
	"runtime"
	"time"
)

func (a *Agent) StartServerMonitor() {
	for {
		a.updateStatus()
		time.Sleep(3 * time.Second)
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
