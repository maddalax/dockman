package router

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"log/slog"
	"net/http"
	"paas/docker"
	"paas/must"
	"paas/resources"
	"time"
)

func StartProxy(locator *service.Locator) {
	service.Set(locator, service.Singleton, func() *multiproxy.LoadBalancer {
		return multiproxy.CreateLoadBalancer()
	})

	lb := service.Get[multiproxy.LoadBalancer](locator)

	go func() {
		for {
			lb.PrintMetrics()
			time.Sleep(5 * time.Second)
		}
	}()

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
			slog.Error("Failed to start server", slog.String("error", err.Error()))
			panic(err)
		}
	}()
}

func ReloadConfig(locator *service.Locator) {
	lb := service.Get[multiproxy.LoadBalancer](locator)
	slog.Info("Reloading reverse proxy upstream config")
	config := loadConfig(locator)
	lb.SetUpstreams(config.Upstreams)
}

func loadConfig(locator *service.Locator) Config {
	var upstreams []*multiproxy.Upstream
	table, err := GetRouteTable(locator)
	if err != nil {
		return Config{
			Upstreams: upstreams,
			Matcher:   &Matcher{},
		}
	}

	matcher := &Matcher{}

	for _, block := range table {

		resource, err := resources.Get(locator, block.ResourceId)
		if err != nil {
			slog.Error("Failed to get resource", slog.String("resourceId", block.ResourceId), slog.String("error", err.Error()))
			continue
		}
		dockerClient, err := docker.Connect()
		if err != nil {
			slog.Error("Failed to connect to docker", slog.String("error", err.Error()))
			continue
		}
		container, err := dockerClient.GetContainer(resource)
		if err != nil {
			slog.Error("Failed to get container", slog.String("error", err.Error()))
			continue
		}

		for port, binding := range container.NetworkSettings.Ports {
			if port.Proto() == "tcp" {
				for _, portBinding := range binding {

					upstream := &multiproxy.Upstream{
						Url: must.Url(fmt.Sprintf("http://%s:%s", portBinding.HostIP, portBinding.HostPort)),
						MatchesFunc: func(req *http.Request, match *multiproxy.Match) bool {
							return matcher.Matches(req)
						},

						// really doesn't matter since we are overriding the MatchesFunc
						Matches: []multiproxy.Match{
							{
								Host: "*",
								Path: "*",
							},
						},
					}

					matcher.AddUpstream(upstream, &block)
					upstreams = append(upstreams, upstream)
				}
			}
		}
	}

	return Config{
		Upstreams: upstreams,
		Matcher:   matcher,
	}
}
