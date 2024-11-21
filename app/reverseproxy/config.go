package reverseproxy

import (
	"dockside/app"
	"dockside/app/logger"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
)

func ReloadConfig(locator *service.Locator) {
	proxy := GetInstance(locator)
	logger.Info("Reloading reverse proxy upstream config")
	config := loadConfig(locator)
	proxy.lb.SetUpstreams(h.Map(config.Upstreams, func(u *UpstreamWithResource) *multiproxy.Upstream {
		return u.Upstream
	}))
	proxy.config = config
}

func loadConfig(locator *service.Locator) *Config {
	matcher := &Matcher{}
	builder := NewConfigBuilder(locator, matcher)
	table, err := GetRouteTable(locator)

	if err != nil {
		return builder.Build()
	}

	for _, block := range table {

		resource, err := app.ResourceGet(locator, block.ResourceId)
		if err != nil {
			logger.ErrorWithFields("Failed to to get resource", err, map[string]any{
				"resourceId": block.ResourceId,
			})
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
