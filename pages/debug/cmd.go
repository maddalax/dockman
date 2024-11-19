package debug

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"paas/internal"
	"time"
)

func CmdPage(ctx *h.RequestContext) *h.Page {
	a := internal.AgentFromLocator(ctx.ServiceLocator())

	responses, err := internal.SendCommand[internal.RunResourceResponse](a, internal.SendCommandOpts{
		Command: &internal.RunResourceCommand{
			ResourceId: "e76ea8a4-2ae3-4983-a197-e0ce7d93d1e4",
		},
		Timeout:           time.Second * 10,
		ExpectedResponses: 2,
	})

	if err != nil {
		fmt.Printf("Failed to send command: %s\n", err.Error())
	} else {
		for i, response := range responses {
			fmt.Printf("Response %d: %s\n", i, response.Response.Message)
		}
	}

	return h.EmptyPage()
}
