package main

import (
	"dockside/__htmgo"
	"dockside/app"
	"dockside/app/reverseproxy"
	"fmt"
	"github.com/maddalax/htmgo/extensions/websocket"
	ws2 "github.com/maddalax/htmgo/extensions/websocket/opts"
	"github.com/maddalax/htmgo/extensions/websocket/session"
	"github.com/maddalax/htmgo/framework/config"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"io/fs"
	"net/http"
)

import _ "net/http/pprof"

func main() {

	locator := service.NewLocator()
	cfg := config.Get()
	agent := app.NewAgent(locator)

	service.Set[app.BuilderRegistry](locator, service.Singleton, func() *app.BuilderRegistry {
		return app.NewBuilderRegistry()
	})

	_, err := app.StartNatsServer()

	if err != nil {
		panic(err)
	}

	agent.Setup()

	//go agent.Run()

	reverseproxy.StartProxy(locator)

	m := app.NewMonitor(locator)
	service.Set(locator, service.Singleton, func() *app.ResourceMonitor {
		return m
	})
	m.Start()

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	// TODO remove
	app.RunSandbox(locator)

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
