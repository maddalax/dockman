package deployment

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/domain"
	"paas/pages/resource/resourceui"
)

func Environment(ctx *h.RequestContext) *h.Page {
	return resourceui.Page(ctx, func(resource *domain.Resource) *h.Element {
		return h.Div(
			h.Pf("Environment"),
		)
	})
}
