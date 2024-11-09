package main

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/config"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"io/fs"
	"net/http"
	"paas/__htmgo"
	"paas/nats"
)

func main() {
	locator := service.NewLocator()
	cfg := config.Get()

	_, err := nats.StartServer()

	if err != nil {
		panic(err)
	}

	client, err := nats.Connect(nats.Options{
		Port: 4222,
	})
	if err != nil {
		panic(err)
	}
	bucket, err := client.GetBucket("test")
	if err != nil {
		panic(err)
	}

	d, err := bucket.Get("test")

	if d != nil {
		fmt.Println(d.Value())
	}

	_, err = bucket.Put("test", []byte("hello"))
	if err != nil {
		panic(err)
	}

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
