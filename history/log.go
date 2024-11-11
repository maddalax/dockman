package history

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/kv"
	"paas/kv/subject"
	"time"
)

func LogChange(locator *service.Locator, subject subject.Subject, data map[string]any) {
	client := service.Get[kv.Client](locator)
	client.CreateHistoryStream()
	data["created_at"] = time.Now().Format(time.Stamp)
	data["subject"] = subject
	client.Publish(subject, kv.MustSerialize(data))
}
