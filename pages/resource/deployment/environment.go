package deployment

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/app"
	"paas/pages/resource/resourceui"
)

func Environment(ctx *h.RequestContext) *h.Page {
	return resourceui.Page(ctx, func(resource *app.Resource) *h.Element {
		return h.Div(
			h.Pf("Environment"),
		)
	})
}
