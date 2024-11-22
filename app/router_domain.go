package app

import (
	"github.com/gobwas/glob"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/maddalax/multiproxy"
	"sync/atomic"
)

type ReverseProxy struct {
	lb            *multiproxy.LoadBalancer[UpstreamMeta]
	locator       *service.Locator
	totalRequests atomic.Int64
}

type RouteBlock struct {
	Hostname          string
	Path              string
	ResourceId        string
	PathMatchModifier string
}

type UpstreamMeta struct {
	Resource     *Resource
	Server       *Server
	Block        *RouteBlock
	GlobPatterns map[string]glob.Glob
}

type CustomUpstream = multiproxy.Upstream[UpstreamMeta]
