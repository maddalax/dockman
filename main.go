package main

import (
	"dockman/__htmgo"
	"dockman/app"
	"dockman/middleware"
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
	registry := app.CreateServiceRegistry(locator)

	app.MustStartNats()

	registry.RegisterStartupServices()

	// Need to register these to be able to send commands even
	// if this process is not running as an agent
	registry.GetAgent().RegisterGobTypes()

	// Setup the reverse proxy
	registry.GetReverseProxy().Setup()

	go registry.GetResourceMonitor().Start()
	go registry.GetReverseProxy().Start()
	go registry.GetJobRunner().Start()

	// TODO remove
	app.RunSandbox(locator)

	h.Start(h.AppOpts{
		ServiceLocator: locator,
		LiveReload:     true,
		Register: func(a *h.App) {
			a.Use(func(ctx *h.RequestContext) {
				session.CreateSession(ctx)
			})

			middleware.UseLoginRequiredMiddleware(a.Router)

			websocket.EnableExtension(a, ws2.ExtensionOpts{
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

			cfg := config.Get()
			// change this in htmgo.yml (public_asset_path)
			a.Router.Handle(fmt.Sprintf("%s/*", cfg.PublicAssetPath),
				http.StripPrefix(cfg.PublicAssetPath, http.FileServerFS(sub)))

			__htmgo.Register(a.Router)
		},
	})
}
