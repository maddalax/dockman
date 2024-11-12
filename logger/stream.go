package logger

import (
	"context"
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"io"
	"paas/docker"
	"paas/kv"
	"paas/kv/subject"
	"paas/resources"
)

type StreamLogsOptions struct {
	Stdout io.Writer
}

func StreamLogs(locator *service.Locator, context context.Context, resource *resources.Resource, cb func(msg *nats.Msg)) *kv.WriterSubscriber {
	natsClient := kv.GetClientFromLocator(locator)
	writer := natsClient.CreateEphemeralWriterSubscriber(context, subject.RunLogsForResource(resource.Id))

	go func() {
		for {
			select {
			case <-context.Done():
				return
			case msg := <-writer.Subscriber:
				cb(msg)
			}
		}
	}()

	switch resource.BuildMeta.(type) {
	case *resources.DockerBuildMeta:
		streamDockerLogs(resource, writer)
	}

	return writer
}

func streamDockerLogs(resource *resources.Resource, writer *kv.WriterSubscriber) {
	client, err := docker.Connect()
	if err != nil {
		writer.Writer.Write([]byte(err.Error()))
		return
	}
	containerId := fmt.Sprintf("%s-%s-container", resource.Name, resource.Id)
	err = client.StreamLogs(containerId, docker.StreamLogsOptions{
		Stdout: writer.Writer,
	})
	if err != nil {
		writer.Writer.Write([]byte(err.Error()))
	}
}
