package app

import (
	"context"
	"github.com/maddalax/htmgo/extensions/websocket/session"
	"paas/app/logger"
)

func DisposeOnCancel(ctx context.Context, dispose func()) {
	go func() {
		<-ctx.Done()
		socketId := session.Id("")
		if v := ctx.Value("socketId"); v != nil {
			socketId = v.(session.Id)
		}
		logger.DebugWithFields("Disposing due to context cancel", map[string]any{
			"socketId": string(socketId),
		})
		dispose()
	}()
}
