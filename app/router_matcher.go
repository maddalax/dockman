package app

import (
	"dockside/app/logger"
	"github.com/gobwas/glob"
	"github.com/maddalax/multiproxy"
	"net/http"
	"strings"
)

type UpstreamWithBlock struct {
	Upstream     *multiproxy.Upstream
	Block        *RouteBlock
	GlobPatterns map[string]glob.Glob
}

func (u *UpstreamWithBlock) Compile() bool {
	if u.GlobPatterns == nil {
		u.GlobPatterns = make(map[string]glob.Glob)
	}
	g, err := glob.Compile(u.Block.Path)
	if err != nil {
		logger.Error("Failed to compile glob", err)
		return false
	}
	u.GlobPatterns[u.Block.Path] = g
	return true
}

func (u *UpstreamWithBlock) Matches(req *http.Request) bool {
	block := u.Block

	if req.Host != block.Hostname {
		return false
	}

	if block.PathMatchModifier == "" {
		block.PathMatchModifier = "starts-with"
	}

	if block.Path == "" {
		return true
	}

	path := req.URL.Path

	switch block.PathMatchModifier {
	case "starts-with":
		return strings.HasPrefix(path, block.Path)
	case "not-starts-with":
		return !strings.HasPrefix(path, block.Path)
	case "not-equals":
		return path != block.Path
	case "not-ends-with":
		return !strings.HasSuffix(path, block.Path)
	case "ends-with":
		return strings.HasSuffix(path, block.Path)
	case "contains":
		return strings.Contains(path, block.Path)
	case "is":
		return path == block.Path
	case "glob":
		g := u.GlobPatterns[block.Path]
		if g != nil {
			return g.Match(path)
		}
		return false
	default:
		// should never happen
		return false
	}
}

type Matcher struct {
	Upstreams []*UpstreamWithBlock
}

func (m *Matcher) AddUpstream(upstream *multiproxy.Upstream, block *RouteBlock) {
	upb := &UpstreamWithBlock{Upstream: upstream, Block: block}
	upb.Compile()
	m.Upstreams = append(m.Upstreams, upb)
}

func (m *Matcher) Matches(req *http.Request) bool {
	for _, upb := range m.Upstreams {
		if upb.Matches(req) {
			return true
		}
	}
	return false
}
