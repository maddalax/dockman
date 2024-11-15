package deployment

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/pages/resource/resourceui"
	"paas/resources"
)

func Environment(ctx *h.RequestContext) *h.Page {
	return resourceui.Page(ctx, func(resource *resources.Resource) *h.Element {
		return h.Div(
			h.Pf("Environment"),
		)
	})
}
