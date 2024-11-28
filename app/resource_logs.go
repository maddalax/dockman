package app

import (
	"context"
	"dockman/app/subject"
	"dockman/app/util/json2"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
)

func StreamResourceLogs(locator *service.Locator, context context.Context, resource *Resource, cb func(log *DockerLog)) error {
	kv := KvFromLocator(locator)
	subjectName := subject.RunLogsForResource(resource.Id)
	streamName := kv.RunLogStreamName(resource.Id)
	streamInfo, err := kv.js.StreamInfo(streamName)
	if err != nil {
		return err
	}

	opts := []nats.SubOpt{
		nats.StartSequence(streamInfo.State.LastSeq - 100),
	}

	_, err = kv.SubscribeStream(context, subjectName, opts, func(msg *nats.Msg) {
		log, err := json2.Deserialize[DockerLog](msg.Data)
		if err == nil {
			cb(log)
		}
	})
	return err
}
