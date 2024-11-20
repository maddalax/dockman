package app

import (
	"fmt"
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
			fmt.Printf("Error attaching server %s to resource: %s\n", server.Id, err.Error())
		}
	}

}
