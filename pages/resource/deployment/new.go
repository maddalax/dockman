package deployment

import (
	"dockman/app"
	"dockman/app/urls"
	"dockman/pages"
	"github.com/google/uuid"
	"github.com/maddalax/htmgo/framework/h"
	"time"
)

func StartNewDeployment(ctx *h.RequestContext) *h.Page {

	resourceId := ctx.QueryParam("resourceId")
	buildId := ctx.QueryParam("buildId")
	isExistingBuild := buildId != ""

	if resourceId == "" {
		ctx.Redirect("/", 302)
		return h.EmptyPage()
	}

	if buildId == "" {
		buildId = uuid.NewString()
	}

	resource, err := app.ResourceGet(ctx.ServiceLocator(), resourceId)

	// todo better error handling
	if err != nil {
		return pages.SidebarPage(
			ctx,
			h.Div(
				h.Pf("failed to find resource"),
			),
		)
	}

	b := app.NewResourceBuilder(ctx.ServiceLocator(), resource, buildId, "Manual (User Requested)")

	if isExistingBuild {
		// starting a new build, clear any previous logs for this build
		b.ClearLogs()
	}

	// waiting 2 seconds so they can see the build log starting
	err = b.StartBuildAsync(time.Second * 2)

	// todo better error handling
	if err != nil {
		return pages.SidebarPage(
			ctx,
			h.Div(
				h.Pf("failed to start build"),
			),
		)
	}

	ctx.Redirect(urls.ResourceDeploymentLogUrl(resourceId, buildId), 302)
	return h.EmptyPage()
}
