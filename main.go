package main

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/config"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"io/fs"
	"net/http"
	"paas/__htmgo"
	"paas/caddy"
	"paas/kv"
)

func main() {
	locator := service.NewLocator()
	cfg := config.Get()

	_, err := kv.StartServer()

	if err != nil {
		panic(err)
	}

	natsClient, err := kv.Connect(kv.Options{
		Port: 4222,
	})
	if err != nil {
		panic(err)
	}

	natsClient.CreateBuildLogStream()

	natsClient.SubscribeAndReplayAll(kv.BuildLogStreamSubject, func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
	})

	go caddy.Run()

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

	h.Start(h.AppOpts{
		ServiceLocator: locator,
		LiveReload:     true,
		Register: func(app *h.App) {
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
