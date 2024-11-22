package app

import (
	"dockside/app/logger"
	"github.com/maddalax/htmgo/framework/service"
)

// UpstreamPortMonitor the resources ports may change if the container is restarted,
// which means the upstreams that the router is aware of needs to be updated. This function
// will monitor the ports of the upstreams and update the router when they change.

func (r *ReverseProxy) UpstreamPortMonitor(locator *service.Locator) {
	loadConfig(locator)
	// if the old lastConfig has a port difference with the new lastConfig, reload the lastConfig
	if r.HasPortDifference() {
		logger.InfoWithFields("Reloading reverse proxy upstream lastConfig", map[string]any{
			"upstreams": len(r.lb.GetStagedUpstreams()),
		})
		r.lb.ApplyStagedUpstreams()
	}
}
