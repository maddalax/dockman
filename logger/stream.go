package logger

import (
	"context"
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"io"
	"paas/docker"
	"paas/domain"
	"paas/kv"
	"paas/kv/subject"
	"paas/monitor"
	"time"
)

type StreamLogsOptions struct {
	Stdout io.WriteCloser
	Since  time.Time
}

func StreamLogs(locator *service.Locator, context context.Context, resource *domain.Resource, cb func(msg *nats.Msg)) {
	doStream(locator, context, resource, cb, time.Time{})
}

func doStream(locator *service.Locator, context context.Context, resource *domain.Resource, cb func(msg *nats.Msg), lastMessageTime time.Time) {
	natsClient := kv.GetClientFromLocator(locator)
	restartStream := false

	writer := natsClient.CreateEphemeralWriterSubscriber(context, subject.RunLogsForResource(resource.Id), kv.CreateOptions{
		BeforeWrite: func(data string) bool {
			lastMessageTime = time.Now()
			return true
		},
	})

	m := service.Get[monitor.Monitor](locator)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	streaming := true

	var startStreaming = func() {
		streaming = true
		switch resource.BuildMeta.(type) {
		case *domain.DockerBuildMeta:
			// this is blocking, so if this stops then we know streaming stopped
			streamDockerLogs(resource, StreamLogsOptions{
				Stdout: writer.Writer,
				Since:  lastMessageTime,
			})
			fmt.Printf("streaming stopped for %s\n", resource.Id)
			streaming = false
		}
	}

	go startStreaming()

	for {
		if restartStream {
			break
		}
		select {
		case <-context.Done():
			return
		case msg := <-writer.Subscriber:
			cb(msg)
		case <-ticker.C:
			if streaming {
				continue
			}
			// streaming stopped, lets see if we need to re-connect it
			status := m.GetRunStatus(resource)
			fmt.Printf("streaming is stopped, checking run status %s\n", resource.Id)
			if status == domain.RunStatusRunning {
				fmt.Printf("container is running, starting streaming again %s\n", resource.Id)
				restartStream = true
				break
			} else {
				fmt.Printf("container is not running, do nothing %s\n", resource.Id)
				// container is not running, do nothing, we'll check again in 3s if it's running
				continue
			}
		}
	}

	fmt.Print("breaking...")
	if restartStream {
		writer = nil
		fmt.Print("restarting stream")
		doStream(locator, context, resource, cb, lastMessageTime)
	}
}

func streamDockerLogs(resource *domain.Resource, opts StreamLogsOptions) {
	client, err := docker.Connect()
	if err != nil {
		opts.Stdout.Write([]byte(err.Error()))
		return
	}
	containerId := fmt.Sprintf("%s-%s-container", resource.Name, resource.Id)
	err = client.StreamLogs(containerId, docker.StreamLogsOptions{
		Stdout: opts.Stdout,
		Since:  opts.Since,
	})
	if err != nil {
		opts.Stdout.Write([]byte(err.Error()))
	}
}
