package app

import (
	"dockside/app/logger"
	"github.com/maddalax/htmgo/framework/service"
)

// UpstreamPortMonitor the resources ports may change if the container is restarted,
// which means the upstreams that the router is aware of needs to be updated. This function
// will monitor the ports of the upstreams and update the router when they change.

func (r *ReverseProxy) UpstreamPortMonitor(locator *service.Locator) {
	newConfig := loadConfig(locator)
	// if the old config has a port difference with the new config, reload the config
	if r.config.HasPortDifference(newConfig) {
		logger.Info("Upstream config difference detected, reloading config")
		ReloadConfig(locator)
	}
}
