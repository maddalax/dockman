package app

import (
	"context"
	"dockside/app/subject"
	"dockside/app/util/json2"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
)

func StreamResourceLogs(locator *service.Locator, context context.Context, resource *Resource, cb func(log *DockerLog)) error {
	streamName := subject.RunLogsForResource(resource.Id)
	kv := KvFromLocator(locator)
	_, err := kv.SubscribeStreamAndReplayAll(context, streamName, func(msg *nats.Msg) {
		log, err := json2.Deserialize[DockerLog](msg.Data)
		if err == nil {
			cb(log)
		}
	})
	return err
}
