package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"log/slog"
	"net/http"
	"time"
)

func StartProxy(locator *service.Locator) {
	service.Set(locator, service.Singleton, func() *multiproxy.LoadBalancer {
		return multiproxy.CreateLoadBalancer()
	})

	lb := service.Get[multiproxy.LoadBalancer](locator)

	if h.IsDevelopment() {
		go func() {
			for {
				lb.PrintMetrics()
				time.Sleep(5 * time.Second)
			}
		}()
	}

	config := loadConfig(locator)
	lb.SetUpstreams(config.Upstreams)

	handler := multiproxy.NewReverseProxyHandler(lb)

	router := chi.NewRouter()
	router.HandleFunc("/*", handler)

	server := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("Failed to start reverse proxy server", slog.String("error", err.Error()))
			panic(err)
		}
	}()
}
