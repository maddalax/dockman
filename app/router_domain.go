package app

import (
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
)

type ReverseProxy struct {
	lb      *multiproxy.LoadBalancer
	config  *Config
	locator *service.Locator
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
	Upstream *multiproxy.Upstream
	Resource *Resource
	Server   *Server
}
