package wsutil

import (
	"context"
	"github.com/maddalax/htmgo/extensions/websocket/session"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"log/slog"
	"time"
)

type OnceAliveOpts struct {
	OnClose func()
}

// OnceWithAliveContext runs a handler with a context that will be cancelled if the websocket connection is lost
// the callback is run once the websocket connection is established
func OnceWithAliveContext(ctx *h.RequestContext, handler func(context.Context)) {
	cc := WithAliveContext(ctx)
	ws.Once(ctx, func() {
		handler(cc)
	})
}

// WithAliveContext creates a context that will be cancelled if the websocket connection is lost
func WithAliveContext(ctx *h.RequestContext) context.Context {
	ccv := context.WithValue(context.Background(), "socketId", session.GetSessionId(ctx))
	cc, cancel := context.WithCancel(ccv)
	socketId := session.GetSessionId(ctx)
	manager := ws.ManagerFromCtx(ctx)

	ws.Once(ctx, func() {
		go func() {
			for {
				if manager.Get(string(socketId)) == nil {
					cancel()
					return
				}
				slog.Debug("socket is alive, waiting for close", slog.String("socketId", string(socketId)))
				time.Sleep(time.Second)
			}
		}()
	})

	return cc
}
