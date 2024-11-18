package router

import (
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"log/slog"
	"paas/resources"
)

func ReloadConfig(locator *service.Locator) {
	lb := service.Get[multiproxy.LoadBalancer](locator)
	slog.Info("Reloading reverse proxy upstream config")
	config := loadConfig(locator)
	lb.SetUpstreams(config.Upstreams)
}

func loadConfig(locator *service.Locator) *Config {
	matcher := &Matcher{}
	builder := NewConfigBuilder(matcher)
	table, err := GetRouteTable(locator)

	if err != nil {
		return builder.Build()
	}

	for _, block := range table {

		resource, err := resources.Get(locator, block.ResourceId)
		if err != nil {
			slog.Error("Failed to get resource", slog.String("resourceId", block.ResourceId), slog.String("error", err.Error()))
			continue
		}

		err = builder.Append(resource, &block)

		if err != nil {
			panic(err)
		}
	}

	return builder.Build()
}
