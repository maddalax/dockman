package main

import (
	"fmt"
	"github.com/maddalax/htmgo/extensions/websocket"
	ws2 "github.com/maddalax/htmgo/extensions/websocket/opts"
	"github.com/maddalax/htmgo/extensions/websocket/session"
	"github.com/maddalax/htmgo/framework/config"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"io/fs"
	"net/http"
	"paas/__htmgo"
	"paas/internal"
	"paas/internal/reverseproxy"
)

import _ "net/http/pprof"

func main() {
	locator := service.NewLocator()
	cfg := config.Get()

	service.Set[internal.BuilderRegistry](locator, service.Singleton, func() *internal.BuilderRegistry {
		return internal.NewBuilderRegistry()
	})

	service.Set[internal.KvClient](locator, service.Singleton, func() *internal.KvClient {
		client, err := internal.NatsConnect(internal.NatsConnectOptions{
			Port: 4222,
		})
		if err != nil {
			panic(err)
		}
		return client
	})

	_, err := internal.StartNatsServer()

	if err != nil {
		panic(err)
	}

	reverseproxy.StartProxy(locator)

	m := internal.NewMonitor(locator)
	service.Set(locator, service.Singleton, func() *internal.ResourceMonitor {
		return m
	})

	go m.StartRunStatusMonitor()

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	h.Start(h.AppOpts{
		ServiceLocator: locator,
		LiveReload:     true,
		Register: func(app *h.App) {

			app.Use(func(ctx *h.RequestContext) {
				session.CreateSession(ctx)
			})

			websocket.EnableExtension(app, ws2.ExtensionOpts{
				WsPath: "/ws",
				RoomName: func(ctx *h.RequestContext) string {
					return "all"
				},
				SessionId: func(ctx *h.RequestContext) string {
					return ctx.QueryParam("sessionId")
				},
			})

			sub, err := fs.Sub(GetStaticAssets(), "assets/dist")

			if err != nil {
				panic(err)
			}

			http.FileServerFS(sub)

			// change this in htmgo.yml (public_asset_path)
			app.Router.Handle(fmt.Sprintf("%s/*", cfg.PublicAssetPath),
				http.StripPrefix(cfg.PublicAssetPath, http.FileServerFS(sub)))

			__htmgo.Register(app.Router)
		},
	})
}
