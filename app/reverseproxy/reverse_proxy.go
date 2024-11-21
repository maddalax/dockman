package reverseproxy

import (
	"dockside/app/logger"
	"github.com/go-chi/chi/v5"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"net/http"
)

func GetInstance(locator *service.Locator) *ReverseProxy {
	return service.Get[ReverseProxy](locator)
}

func createReverseProxy(locator *service.Locator) *ReverseProxy {
	lb := multiproxy.CreateLoadBalancer()
	config := loadConfig(locator)
	lb.SetUpstreams(h.Map(config.Upstreams, func(u *UpstreamWithResource) *multiproxy.Upstream {
		return u.Upstream
	}))
	return &ReverseProxy{
		lb:     lb,
		config: config,
	}
}

func StartProxy(locator *service.Locator) {
	created := createReverseProxy(locator)

	service.Set(locator, service.Singleton, func() *ReverseProxy {
		return created
	})

	proxy := GetInstance(locator)

	// Start the upstream port monitor to detect changes in the upstreams
	go proxy.StartUpstreamPortMonitor(locator)

	handler := multiproxy.NewReverseProxyHandler(proxy.lb)

	router := chi.NewRouter()
	router.HandleFunc("/*", handler)

	server := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("Failed to start reverse proxy server", err)
			panic(err)
		}
	}()
}
