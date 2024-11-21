package app

import (
	"dockside/app/logger"
	"dockside/app/subject"
	"dockside/app/util/json2"
	"github.com/maddalax/htmgo/framework/service"
	"time"
)

func LogChange(locator *service.Locator, subject subject.Subject, data map[string]any) {
	client := service.Get[KvClient](locator)
	err := client.CreateHistoryStream()
	if err != nil {
		logger.Error("failed to create history stream", err)
		return
	}
	data["created_at"] = time.Now().Format(time.Stamp)
	data["subject"] = subject
	err = client.Publish(subject, json2.SerializeOrEmpty(data))
	if err != nil {
		logger.Error("failed to publish history log", err)
	}
}
