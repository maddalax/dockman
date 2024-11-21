package reverseproxy

import (
	"dockside/app/logger"
	"github.com/maddalax/htmgo/framework/service"
	"time"
)

// StartUpstreamPortMonitor the resources ports may change if the container is restarted,
// which means the upstreams that the router is aware of needs to be updated. This function
// will monitor the ports of the upstreams and update the router when they change.

func (r *ReverseProxy) StartUpstreamPortMonitor(locator *service.Locator) {
	for {
		newConfig := loadConfig(locator)
		// if the old config has a port difference with the new config, reload the config
		if r.config.HasPortDifference(newConfig) {
			logger.Info("Upstream config difference detected, reloading config")
			ReloadConfig(locator)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
