package app

import (
	"context"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"paas/app/subject"
	"paas/app/util/json2"
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
