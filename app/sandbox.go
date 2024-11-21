package app

import (
	"dockside/app/logger"
	"github.com/maddalax/htmgo/framework/service"
)

func RunSandbox(locator *service.Locator) {
	servers, err := ServerList(locator)

	if err != nil {
		return
	}

	for _, server := range servers {
		err := AttachServerToResource(locator, server.Id, "e76ea8a4-2ae3-4983-a197-e0ce7d93d1e4")
		if err != nil {
			logger.ErrorWithFields("Error attaching server to resource", err, map[string]interface{}{
				"server_id": server.Id,
			})
		}
	}
}
