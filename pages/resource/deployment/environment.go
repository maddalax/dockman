package deployment

import (
	"dockside/app"
	"dockside/pages/resource/resourceui"
	"github.com/maddalax/htmgo/framework/h"
)

func Environment(ctx *h.RequestContext) *h.Page {
	return resourceui.Page(ctx, func(resource *app.Resource) *h.Element {
		return h.Div(
			h.Pf("Environment"),
		)
	})
}
