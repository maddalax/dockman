package app

import (
	"encoding/json"
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
		// don't want to log this error since it will be a lot of noise
		return
	}
	if dockerLog.Log != "" && dockerLog.ResourceId != "" {
		a.kv.LogRunMessage(dockerLog.ResourceId, dockerLog.Log)
	}
}
