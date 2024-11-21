package app

import (
	"dockside/app/logger"
	"dockside/app/util/must"
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"github.com/pkg/errors"
	"net/http"
	"sync"
	"time"
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
	servers, err := ResourceGetServers(b.serviceLocator, resource.Id)

	if err != nil {
		return err
	}

	type entry struct {
		resourceId string
		server     *Server
		index      int
	}

	entries := make([]entry, 0)

	for _, server := range servers {
		for i := range resource.InstancesPerServer {
			entries = append(entries, entry{
				resourceId: resource.Id,
				server:     server,
				index:      i,
			})
		}
	}

	wg := sync.WaitGroup{}

	for _, e := range entries {
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch resource.RunType {
			case RunTypeDockerBuild:
				fallthrough
			case RunTypeDockerRegistry:
				err = b.appendDockerUpstreams(resource, e.index, e.server, block)
				if err != nil {
					logger.ErrorWithFields("Failed to append docker upstreams", err, map[string]interface{}{
						"resourceId": resource.Id,
						"serverId":   e.server.Id,
					})
				}
			default:
			}
		}()
	}

	wg.Wait()

	return nil
}

func (b *ConfigBuilder) appendDockerUpstreams(resource *Resource, index int, server *Server, block *RouteBlock) error {

	res, err := SendCommand[GetContainerResponse](b.serviceLocator, server.Id, SendCommandOpts{
		Command: &GetContainerCommand{
			ResourceId:   resource.Id,
			Index:        index,
			ResponseData: &GetContainerResponse{},
		},
		Timeout: time.Second * 5,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to get container for server")
	}

	if res.Response.Error != nil {
		return errors.Wrap(res.Response.Error, "Failed to get container for server")
	}

	if res.Response.Container.Config == nil {
		return errors.Wrap(err, "Failed to get container for server")
	}

	container := res.Response.Container
	hostIp := ""

	if server.RemoteIpAddress != "" {
		hostIp = server.RemoteIpAddress
	}

	// route using local ip first if possible
	if server.LocalIpAddress != "" {
		hostIp = server.LocalIpAddress
	}

	if hostIp == "" {
		return errors.Wrap(err, "Failed to get host ip for server")
	}

	for port, binding := range container.NetworkSettings.Ports {
		if port.Proto() == "tcp" {
			for _, portBinding := range binding {

				upstream := &multiproxy.Upstream{
					Url: must.Url(fmt.Sprintf("http://%s:%s", hostIp, portBinding.HostPort)),
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
