package app

import (
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"log/slog"
)

func ReloadConfig(locator *service.Locator) {
	proxy := GetInstance(locator)
	slog.Info("Reloading reverse proxy upstream config")
	config := loadConfig(locator)
	proxy.lb.SetUpstreams(h.Map(config.Upstreams, func(u *UpstreamWithResource) *multiproxy.Upstream {
		return u.Upstream
	}))
	proxy.config = config
}

func loadConfig(locator *service.Locator) *Config {
	matcher := &Matcher{}
	builder := NewConfigBuilder(matcher)
	table, err := GetRouteTable(locator)

	if err != nil {
		return builder.Build()
	}

	for _, block := range table {

		resource, err := ResourceGet(locator, block.ResourceId)
		if err != nil {
			slog.Error("Failed to get resource", slog.String("resourceId", block.ResourceId), slog.String("error", err.Error()))
			continue
		}

		err = builder.Append(resource, &block)

		if err != nil {
			continue
		}
	}

	return builder.Build()
}

func (c *Config) HasPortDifference(old *Config) bool {
	if len(old.Upstreams) != len(c.Upstreams) {
		return true
	}

	for i, u := range c.Upstreams {
		if u.Upstream.Url.Port() != old.Upstreams[i].Upstream.Url.Port() {
			return true
		}
	}

	return false
}
