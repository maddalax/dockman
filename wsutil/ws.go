package wsutil

import (
	"context"
	"github.com/maddalax/htmgo/extensions/websocket/session"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
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

	listener := make(chan ws.SocketEvent)
	manager := ws.ManagerFromCtx(ctx)

	go func() {
		for {
			select {
			case event := <-listener:
				if event.Type == ws.DisconnectedEvent && event.SessionId == string(socketId) {
					manager.RemoveListener(listener)
					cancel()
					return
				}
			}
		}
	}()

	manager.Listen(listener)

	return cc
}
