package app

import (
	"dockside/app/logger"
	"errors"
	"github.com/gobwas/glob"
	"net/http"
	"strings"
)

func UpstreamMatches(up *CustomUpstream, req *http.Request) bool {
	block := up.Metadata.Block

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
		g := up.Metadata.GlobPatterns[block.Path]
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
}

func (m *Matcher) CompileUpstream(u *CustomUpstream) {
	if u.Metadata.GlobPatterns == nil {
		u.Metadata.GlobPatterns = make(map[string]glob.Glob)
	}
	if u.Metadata.Block == nil {
		logger.ErrorWithFields("Block is nil", errors.New("upstream should have block"), map[string]interface{}{
			"upstream": u.Id,
		})
		return
	}
	g, err := glob.Compile(u.Metadata.Block.Path)
	if err != nil {
		logger.Error("Failed to compile glob", err)
		return
	}
	u.Metadata.GlobPatterns[u.Metadata.Block.Path] = g
}
