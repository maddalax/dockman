package deployment

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/domain"
	"paas/pages/resource/resourceui"
	"paas/ui"
)

func BuidLog(ctx *h.RequestContext) *h.Page {
	buildId := ctx.QueryParam("buildId")
	return resourceui.Page(
		ctx, func(resource *domain.Resource) *h.Element {
			return h.Div(
				ui.DockerBuildLogs(ctx, resource, buildId),
			)
		},
	)
}
