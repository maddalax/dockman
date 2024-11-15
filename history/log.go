package history

import (
	"github.com/maddalax/htmgo/framework/service"
	"log/slog"
	"paas/kv"
	"paas/kv/subject"
	"time"
)

func LogChange(locator *service.Locator, subject subject.Subject, data map[string]any) {
	client := service.Get[kv.Client](locator)
	err := client.CreateHistoryStream()
	if err != nil {
		slog.Error("failed to create history stream: %v", err)
		return
	}
	data["created_at"] = time.Now().Format(time.Stamp)
	data["subject"] = subject
	err = client.Publish(subject, kv.MustSerialize(data))
	if err != nil {
		slog.Error("failed to publish history: %v", err)
	}
}
