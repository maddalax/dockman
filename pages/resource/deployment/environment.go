package deployment

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/internal"
	"paas/pages/resource/resourceui"
)

func Environment(ctx *h.RequestContext) *h.Page {
	return resourceui.Page(ctx, func(resource *internal.Resource) *h.Element {
		return h.Div(
			h.Pf("Environment"),
		)
	})
}
