package deployment

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/internal"
	"paas/internal/ui"
	"paas/internal/urls"
	"paas/pages/resource/resourceui"
	"time"
)

func Deployment(ctx *h.RequestContext) *h.Page {
	id := ctx.QueryParam("id")

	resource, err := internal.ResourceGet(ctx.ServiceLocator(), id)

	if err != nil {
		ctx.Redirect("/", 302)
		return h.EmptyPage()
	}

	deployments, err := internal.GetDeployments(ctx.ServiceLocator(), resource.Id)

	if err != nil {
		deployments = []internal.Deployment{}
	}

	return resourceui.Page(ctx, func(resource *internal.Resource) *h.Element {
		return h.Div(
			h.Div(
				h.Class("flex gap-2 items-center"),
				ui.PrimaryButton(ui.ButtonProps{
					Text: "ResourceStart Build",
					Href: urls.ResourceStartDeploymentPath(resource.Id, ""),
				}),
			),
			List(deployments),
		)
	})
}

func List(deployments []internal.Deployment) *h.Element {
	return h.Div(
		h.Class("flex flex-col gap-4 max-w-md"), // Increase gap for better spacing between cards
		h.List(deployments, func(deployment internal.Deployment, index int) *h.Element {
			return h.Div(
				h.Class("bg-white shadow-md rounded-lg overflow-hidden border border-gray-200"), // Card styling
				h.Div(
					h.Class("p-4"), // Padding for card content
					h.Div(
						h.Class("flex justify-between items-center mb-2"),
						h.Div(
							h.Class("flex flex-col"),
							h.Pf(
								"Build ID: %s",
								deployment.BuildId,
							),
							h.Pf(
								"Created At: %s",
								deployment.CreatedAt.Format(time.Stamp),
							),
							h.Pf(
								"Status: %s",
								deployment.Status,
							),
						),
					),
					h.Div(
						h.Class("flex justify-start mt-4"),
						ui.Button(ui.ButtonProps{
							Text: "View Log",
							Href: urls.ResourceDeploymentLogUrl(deployment.ResourceId, deployment.BuildId),
						}),
					),
				),
			)
		}),
	)
}
