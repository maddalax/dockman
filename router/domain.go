package router

import "github.com/maddalax/multiproxy"

type ReverseProxy struct {
	lb     *multiproxy.LoadBalancer
	config *Config
}

type RouteBlock struct {
	Hostname          string
	Path              string
	ResourceId        string
	PathMatchModifier string
}

type Config struct {
	Upstreams []*UpstreamWithResource
	Matcher   *Matcher
}

type UpstreamWithResource struct {
	Upstream   *multiproxy.Upstream
	ResourceId string
}
