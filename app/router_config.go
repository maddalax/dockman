package app

import (
	"dockman/app/logger"
	"github.com/maddalax/htmgo/framework/service"
	"slices"
	"strings"
)

// ReloadConfig force reloads the router configuration
func ReloadConfig(locator *service.Locator) {
	loadConfig(locator)
	lb := GetServiceRegistry(locator).GetReverseProxy().lb
	lb.ApplyStagedUpstreams()
}

// loadConfig calculates the new configuration for the router, but does not apply it,
// it must be applied by calling ApplyStagedUpstreams
func loadConfig(locator *service.Locator) {
	builder := NewConfigBuilder(locator)
	table, err := GetRouteTable(locator)
	if err != nil {
		return
	}

	lb := GetServiceRegistry(locator).GetReverseProxy().lb

	// start the staging process
	lb.ClearStagedUpstreams()

	for _, block := range table {

		resource, err := ResourceGet(locator, block.ResourceId)

		if err != nil {
			logger.ErrorWithFields("Failed to to get resource", err, map[string]any{
				"resourceId": block.ResourceId,
			})
			continue
		}

		err = builder.Append(resource, &block, lb)

		if err != nil {
			continue
		}
	}
}

func (r *ReverseProxy) HasPortDifference() bool {
	current := r.lb.GetUpstreams()
	staged := r.lb.GetStagedUpstreams()

	if len(current) != len(staged) {
		return true
	}

	slices.SortFunc(current, func(a, b *CustomUpstream) int {
		return strings.Compare(a.Id, b.Id)
	})

	slices.SortFunc(staged, func(a, b *CustomUpstream) int {
		return strings.Compare(a.Id, b.Id)
	})

	for i, u := range current {
		if u.Id != staged[i].Id {
			return true
		}
	}

	return false
}
