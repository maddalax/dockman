package app

import (
	"dockside/app/logger"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"slices"
	"strings"
)

func ReloadConfig(locator *service.Locator) {
	proxy := GetServiceRegistry(locator).GetReverseProxy()
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

		resource, err := ResourceGet(locator, block.ResourceId)
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

	// sort them so that the order is consistent when comparing old vs new config
	slices.SortFunc(builder.upstreams, func(a, b *UpstreamWithResource) int {
		key := fmt.Sprintf("%s%s", a.Upstream.Url.Host, a.Upstream.Url.Port())
		key2 := fmt.Sprintf("%s%s", b.Upstream.Url.Host, b.Upstream.Url.Port())
		return strings.Compare(key, key2)
	})

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
