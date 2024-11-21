package deployment

import (
	"dockside/app"
	"dockside/app/ui"
	"dockside/pages/resource/resourceui"
	"github.com/maddalax/htmgo/framework/h"
)

func BuidLog(ctx *h.RequestContext) *h.Page {
	buildId := ctx.QueryParam("buildId")
	return resourceui.Page(
		ctx, func(resource *app.Resource) *h.Element {
			return h.Div(
				ui.DockerBuildLogs(ctx, resource, buildId),
			)
		},
	)
}
