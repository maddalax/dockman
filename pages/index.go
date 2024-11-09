package pages

import (
	"fmt"
	"github.com/maddalax/htmgo/extensions/websocket/session"
	"github.com/maddalax/htmgo/framework/h"
	"paas/ui"
)

func IndexPage(ctx *h.RequestContext) *h.Page {
	sessionId := session.GetSessionId(ctx)

	return RootPage(
		h.Div(
			h.Attribute("ws-connect", fmt.Sprintf("/ws?sessionId=%s", sessionId)),
			h.TriggerChildren(),
			h.Class("flex flex-col gap-4 items-center pt-24 min-h-screen bg-neutral-100"),
			ui.DockerBuildTest(ctx),
		),
	)
}
