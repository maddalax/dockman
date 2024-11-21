package app

import (
	"dockside/app/logger"
	"github.com/go-chi/chi/v5"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"net/http"
	"time"
)

func CreateReverseProxy(locator *service.Locator) *ReverseProxy {
	lb := multiproxy.CreateLoadBalancer()
	return &ReverseProxy{
		lb:      lb,
		config:  &Config{},
		locator: locator,
	}
}

func (r *ReverseProxy) Setup() {
	registry := GetServiceRegistry(r.locator)
	// Start the upstream port monitor to detect changes in the upstreams
	registry.GetJobRunner().Add("ReverseProxyCheckUpstreamPorts", time.Second*2, func() {
		r.UpstreamPortMonitor(r.locator)
	})
}

func (r *ReverseProxy) GetConfig() *Config {
	return r.config
}

func (r *ReverseProxy) Start() {
	ReloadConfig(r.locator)

	handler := multiproxy.NewReverseProxyHandler(r.lb)

	router := chi.NewRouter()
	router.HandleFunc("/*", handler)

	server := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start reverse proxy server", err)
		panic(err)
	}
}
