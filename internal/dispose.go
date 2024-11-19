package internal

import (
	"context"
	"github.com/maddalax/htmgo/extensions/websocket/session"
	"log/slog"
)

func DisposeOnCancel(ctx context.Context, dispose func()) {
	go func() {
		<-ctx.Done()
		socketId := session.Id("")
		if v := ctx.Value("socketId"); v != nil {
			socketId = v.(session.Id)
		}
		slog.Debug("Disposing due to context cancel", slog.String("socketId", string(socketId)))
		dispose()
	}()
}
