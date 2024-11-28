package app

import (
	"dockside/app/logger"
	"github.com/go-chi/chi/v5"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"net/http"
	"sync/atomic"
	"time"
)

func CreateReverseProxy(locator *service.Locator) *ReverseProxy {
	lb := multiproxy.CreateLoadBalancer[UpstreamMeta]()
	return &ReverseProxy{
		lb:            lb,
		locator:       locator,
		totalRequests: atomic.Int64{},
	}
}

func (r *ReverseProxy) Setup() {
	registry := GetServiceRegistry(r.locator)
	source := "dockside"
	// Start the upstream port monitor to detect changes in the upstreams
	registry.GetJobRunner().Add(source, "ReverseProxyCheckUpstreamPorts", "Checks the ports the running containers are on for each connected server, so the reverse proxy knows where to route to.", time.Second*2, func() {
		r.UpstreamPortMonitor(r.locator)
	})
}

func (r *ReverseProxy) GetUpstreams() []*CustomUpstream {
	return r.lb.GetUpstreams()
}

func (r *ReverseProxy) Start() {
	ReloadConfig(r.locator)

	handler := multiproxy.NewReverseProxyHandler(r.lb)

	router := chi.NewRouter()
	router.HandleFunc("/*", func(writer http.ResponseWriter, request *http.Request) {
		r.totalRequests.Add(1)
		handler(writer, request)
	})

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
