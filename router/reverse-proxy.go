package router

import (
	"github.com/maddalax/htmgo/framework/service"
)

func StartProxy(locator *service.Locator) {
	//table, err := GetRouteTable(locator)
	//if err != nil {
	//	panic(err)
	//}
	//lb := multiproxy.CreateLoadBalancer()
	//for _, block := range table {
	//	resource, err := resources.Get(locator, block.ResourceId)
	//	if err != nil {
	//		slog.Error("Failed to get resource", slog.String("resourceId", block.ResourceId), slog.String("error", err.Error()))
	//		continue
	//	}
	//	lb.Add(&multiproxy.Upstream{
	//		Url: resource.Url,
	//	})
	//}
}
