package router

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/go-chi/chi/v5"
	"github.com/gobwas/glob"
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
	regex := regexp2.MustCompile(`^(?!\/example).*`, 0)

	if ok, _ := regex.MatchString("example"); ok {
		fmt.Println("Matched")
	}

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

	upstreams := loadUpstreams(locator)
	lb.SetUpstreams(upstreams)

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
	lb.SetUpstreams(loadUpstreams(locator))
}

func loadUpstreams(locator *service.Locator) []*multiproxy.Upstream {
	var upstreams []*multiproxy.Upstream
	table, err := GetRouteTable(locator)
	if err != nil {
		return []*multiproxy.Upstream{}
	}
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
					upstreams = append(upstreams, &multiproxy.Upstream{
						Url: must.Url(fmt.Sprintf("http://%s:%s", portBinding.HostIP, portBinding.HostPort)),
						PathMatchesFunc: func(path string, match *multiproxy.Match) bool {
							if path == match.Path {
								return true
							}
							// todo precompile these
							g, err := glob.Compile(match.Path)
							if err != nil {
								slog.Error("Failed to compile glob", slog.String("error", err.Error()))
								return false
							}
							return g.Match(path)
						},
						Matches: []multiproxy.Match{
							{
								Host: block.Hostname,
								Path: block.Path,
							},
						},
					})
				}
			}
		}
	}
	return upstreams
}
