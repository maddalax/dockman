package app

import (
	"fmt"
	"github.com/maddalax/multiproxy"
	"net/http"
	"paas/must"
)

type ConfigBuilder struct {
	matcher   *Matcher
	upstreams []*UpstreamWithResource
}

func NewConfigBuilder(matcher *Matcher) *ConfigBuilder {
	return &ConfigBuilder{
		matcher:   matcher,
		upstreams: make([]*UpstreamWithResource, 0),
	}
}

func (b *ConfigBuilder) Build() *Config {
	return &Config{
		Upstreams: b.upstreams,
		Matcher:   b.matcher,
	}
}

func (b *ConfigBuilder) Append(resource *Resource, block *RouteBlock) error {
	switch resource.RunType {
	case RunTypeDockerBuild:
		fallthrough
	case RunTypeDockerRegistry:
		for i := range resource.InstancesPerServer {
			err := b.appendDockerUpstreams(resource, i, block)
			if err != nil {
				return err
			}
		}
	default:
	}
	return nil
}

func (b *ConfigBuilder) appendDockerUpstreams(resource *Resource, index int, block *RouteBlock) error {
	dockerClient, err := DockerConnect()
	if err != nil {
		return DockerConnectionError
	}

	container, err := dockerClient.GetContainer(resource, index)
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
				b.upstreams = append(b.upstreams, &UpstreamWithResource{
					Upstream:   upstream,
					ResourceId: resource.Id,
				})
			}
		}
	}

	return nil
}
