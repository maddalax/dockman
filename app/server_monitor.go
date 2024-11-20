package app

import (
	"fmt"
	"log/slog"
	"os"
	"paas/app/util/networkutil"
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
		slog.Error("Failed to update server status", slog.String("error", err.Error()), slog.String("server_id", a.serverId))
	}
}
