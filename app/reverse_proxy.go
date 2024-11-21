package app

import (
	"dockside/app/logger"
	"github.com/go-chi/chi/v5"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"net/http"
)

func CreateReverseProxy(locator *service.Locator) *ReverseProxy {
	lb := multiproxy.CreateLoadBalancer()
	config := loadConfig(locator)
	lb.SetUpstreams(h.Map(config.Upstreams, func(u *UpstreamWithResource) *multiproxy.Upstream {
		return u.Upstream
	}))
	return &ReverseProxy{
		lb:      lb,
		config:  config,
		locator: locator,
	}
}

func (r *ReverseProxy) Start() {
	// Start the upstream port monitor to detect changes in the upstreams
	go r.StartUpstreamPortMonitor(r.locator)

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
