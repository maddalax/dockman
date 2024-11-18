package router

import "github.com/maddalax/multiproxy"

type RouteBlock struct {
	Hostname          string
	Path              string
	ResourceId        string
	PathMatchModifier string
}

type Config struct {
	Upstreams []*multiproxy.Upstream
	Matcher   *Matcher
}
