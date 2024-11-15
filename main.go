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
	"paas/kv"
	"paas/monitor"
	"paas/router"
)

func main() {
	locator := service.NewLocator()
	cfg := config.Get()

	service.Set[kv.Client](locator, service.Singleton, func() *kv.Client {
		client, err := kv.Connect(kv.Options{
			Port: 4222,
		})
		if err != nil {
			panic(err)
		}
		return client
	})

	_, err := kv.StartServer()

	if err != nil {
		panic(err)
	}

	router.StartProxy(locator)

	m := monitor.NewMonitor(locator)
	service.Set(locator, service.Singleton, func() *monitor.Monitor {
		return m
	})

	go m.StartRunStatusMonitor()

	//go func() {
	//	client, err := docker.Connect()
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	outputStream := natsClient.NewNatsWriter(kv.BuildLogStreamSubject)
	//	err = client.Build(outputStream, ".", types.ImageBuildOptions{
	//		Dockerfile: "Dockerfile",
	//	})
	//	if err != nil {
	//		panic(err)
	//	}
	//}()

	//err = resources.Create(locator, resources.CreateOptions{
	//	Name:        "htmgo-" + uuid.NewString()[0:4],
	//	Environment: "dev",
	//	RunType:     resources.RunTypeDockerBuild,
	//	BuildMeta: resources.DockerBuildMeta{
	//		Dockerfile: "Dockerfile",
	//		Tags:       []string{"test", "test2"},
	//	},
	//	Env: []resources.Env{
	//		{
	//			Key:   "TEST",
	//			Value: "test",
	//		},
	//		{
	//			Key:   "TEST2",
	//			Value: "test2",
	//		},
	//	},
	//})
	//
	//if err != nil {
	//	panic(err)
	//}

	h.Start(h.AppOpts{
		ServiceLocator: locator,
		LiveReload:     true,
		Register: func(app *h.App) {

			app.Use(func(ctx *h.RequestContext) {
				host := ctx.Request.Host
				fmt.Printf("Host: %s\n", host)

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
