package app

import (
	"encoding/json"
	"time"
)

type DockerLog struct {
	ContainerId   string    `json:"container_id"`
	Log           string    `json:"log"`
	ContainerName string    `json:"container_name"`
	BuildId       string    `json:"paas.build.id"`
	ResourceId    string    `json:"paas.resource.id"`
	Time          time.Time `json:"timestamp"`
}

func (a *Agent) WriteContainerLog(log string) {
	dockerLog := DockerLog{}
	err := json.Unmarshal([]byte(log), &dockerLog)
	if err != nil {
		// don't want to log this error since it will be a lot of noise
		return
	}
	if dockerLog.Log != "" && dockerLog.ResourceId != "" {
		layout := "2006/01/02 15:04:05"
		parsedTime, err := time.Parse(layout, dockerLog.Log[0:19])
		if err == nil {
			dockerLog.Time = parsedTime
			dockerLog.Log = dockerLog.Log[20:]
		}

		serialized, err := json.Marshal(dockerLog)

		if err != nil {
			return
		}

		a.kv.LogRunMessage(dockerLog.ResourceId, string(serialized))
	}
}
