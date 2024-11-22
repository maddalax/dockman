package app

import (
	"dockside/app/util/must"
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"net/http"
)

type ConfigBuilder struct {
	matcher        *Matcher
	upstreams      []*UpstreamWithResource
	serviceLocator *service.Locator
}

func NewConfigBuilder(locator *service.Locator, matcher *Matcher) *ConfigBuilder {
	return &ConfigBuilder{
		matcher:        matcher,
		upstreams:      make([]*UpstreamWithResource, 0),
		serviceLocator: locator,
	}
}

func (b *ConfigBuilder) Build() *Config {
	return &Config{
		Upstreams: b.upstreams,
		Matcher:   b.matcher,
	}
}

func (b *ConfigBuilder) Append(resource *Resource, block *RouteBlock) error {

	if len(resource.ServerDetails) == 0 {
		return nil
	}

	for _, serverDetail := range resource.ServerDetails {
		if serverDetail.RunStatus == RunStatusNotRunning {
			continue
		}
		server, err := ServerGet(b.serviceLocator, serverDetail.ServerId)
		if err != nil {
			continue
		}

		// skip if server is not accessible
		if !server.IsAccessible() {
			continue
		}

		for _, up := range serverDetail.Upstreams {
			upstream := &multiproxy.Upstream{
				Url: must.Url(fmt.Sprintf("http://%s:%s", up.Host, up.Port)),
				MatchesFunc: func(req *http.Request, match *multiproxy.Match) bool {
					return b.matcher.Matches(req)
				},

				// really doesn't matter since we are overriding the MatchesFunc
				Matches: []multiproxy.Match{
					{
						Host: "*",
						Path: "*",
					},
				},
			}

			b.matcher.AddUpstream(upstream, block)
			b.upstreams = append(b.upstreams, &UpstreamWithResource{
				Upstream: upstream,
				Resource: resource,
				Server:   server,
			})
		}
	}
	return nil
}
