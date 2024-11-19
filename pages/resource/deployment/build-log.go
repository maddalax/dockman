package deployment

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/internal"
	"paas/internal/ui"
	"paas/pages/resource/resourceui"
)

func BuidLog(ctx *h.RequestContext) *h.Page {
	buildId := ctx.QueryParam("buildId")
	return resourceui.Page(
		ctx, func(resource *internal.Resource) *h.Element {
			return h.Div(
				ui.DockerBuildLogs(ctx, resource, buildId),
			)
		},
	)
}
