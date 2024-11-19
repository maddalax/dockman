package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type DockerLog struct {
	ContainerId   string `json:"container_id"`
	Log           string `json:"log"`
	ContainerName string `json:"container_name"`
	BuildId       string `json:"paas.build.id"`
	ResourceId    string `json:"paas.resource.id"`
}

func (a *Agent) WriteContainerLog(log string) {
	dockerLog := DockerLog{}
	err := json.Unmarshal([]byte(log), &dockerLog)
	if err != nil {
		slog.Error("Failed to unmarshal log", slog.String("error", err.Error()), slog.String("log", log))
		return
	}
	if dockerLog.Log != "" && dockerLog.ResourceId != "" {
		err = a.kv.CreateRunLogStream(dockerLog.ResourceId)
		if err != nil {
			fmt.Printf("Failed to create run log stream: %s\n", err.Error())
		}
		a.kv.LogRunMessage(dockerLog.ResourceId, dockerLog.Log)
	}
}
