package router

import (
	"fmt"
	"github.com/maddalax/multiproxy"
	"net/http"
	"paas/docker"
	"paas/domain"
	"paas/must"
)

type ConfigBuilder struct {
	matcher   *Matcher
	upstreams []*multiproxy.Upstream
}

func NewConfigBuilder(matcher *Matcher) *ConfigBuilder {
	return &ConfigBuilder{
		matcher:   matcher,
		upstreams: make([]*multiproxy.Upstream, 0),
	}
}

func (b *ConfigBuilder) Build() *Config {
	return &Config{
		Upstreams: b.upstreams,
		Matcher:   b.matcher,
	}
}

func (b *ConfigBuilder) Append(resource *domain.Resource, block *RouteBlock) error {
	switch resource.RunType {
	case domain.RunTypeDockerBuild:
		fallthrough
	case domain.RunTypeDockerRegistry:
		err := b.appendDockerUpstreams(resource, block)
		if err != nil {
			return err
		}
	default:
	}

	return nil
}

func (b *ConfigBuilder) appendDockerUpstreams(resource *domain.Resource, block *RouteBlock) error {
	dockerClient, err := docker.Connect()
	if err != nil {
		return domain.DockerConnectionError
	}

	container, err := dockerClient.GetContainer(resource)
	if err != nil {
		return err
	}

	for port, binding := range container.NetworkSettings.Ports {
		if port.Proto() == "tcp" {
			for _, portBinding := range binding {

				upstream := &multiproxy.Upstream{
					Url: must.Url(fmt.Sprintf("http://%s:%s", portBinding.HostIP, portBinding.HostPort)),
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
				b.upstreams = append(b.upstreams, upstream)
			}
		}
	}

	return nil
}
