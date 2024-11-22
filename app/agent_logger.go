package app

import (
	"bytes"
	"dockside/app/logger"
	"encoding/json"
	"github.com/maddalax/htmgo/framework/service"
	"time"
)

type DockerLog struct {
	ContainerId   string    `json:"container_id"`
	Log           string    `json:"log"`
	ContainerName string    `json:"container_name"`
	BuildId       string    `json:"dockside.build.id"`
	ResourceId    string    `json:"dockside.resource.id"`
	Time          time.Time `json:"timestamp"`
	HostName      string    `json:"hostname"`
}

type NatsContainerLogWriter struct {
	server  *Server
	locator *service.Locator
	agent   *Agent
}

func NewNatsContainerLogWriter(locator *service.Locator, server *Server, agent *Agent) *NatsContainerLogWriter {
	return &NatsContainerLogWriter{
		server:  server,
		locator: locator,
		agent:   agent,
	}
}

func (w *NatsContainerLogWriter) Write(p []byte) (int, error) {
	lines := bytes.Split(p, []byte("\n"))

	for _, line := range lines {

		if len(line) == 0 {
			continue
		}

		dockerLog := DockerLog{}
		err := json.Unmarshal(line, &dockerLog)

		if err != nil {
			// don't want to log this error since it will be a lot of noise
			logger.ErrorWithFields("Failed to deserialize log message", err, map[string]interface{}{
				"message": string(line),
			})
			return len(p), nil
		}

		if dockerLog.Log != "" && dockerLog.ResourceId != "" {
			layout := "2006/01/02 15:04:05"
			parsedTime, err := time.Parse(layout, dockerLog.Log[0:19])

			if err == nil {
				dockerLog.Time = parsedTime
				dockerLog.Log = dockerLog.Log[20:]
			}

			if w.server != nil {
				dockerLog.HostName = w.server.FormattedName()
			}

			serialized, err := json.Marshal(dockerLog)

			if err != nil {
				logger.ErrorWithFields("Failed to serialize log message", err, map[string]interface{}{
					"message": string(line),
				})
			} else {
				w.agent.registry.KvClient().LogRunMessageBytes(dockerLog.ResourceId, serialized)
			}
		}
	}

	return len(p), nil
}
