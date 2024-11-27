package create

import (
	"dockside/app/ui"
	"github.com/maddalax/htmgo/framework/h"
)

func AdditionalCreateResourceFields(ctx *h.RequestContext) *h.Partial {
	deploymentType := ctx.QueryParam("deployment_type")
	return h.SwapPartial(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-2"),
			h.Id("additional-create-resource-fields"),
			AdditionalFieldsForDeploymentType(ctx, deploymentType),
			EnvironmentVariables(ctx),
			h.Div(
				h.Id("submit-error"),
			),
			h.Div(
				ui.SubmitButton(ui.ButtonProps{
					FullWidth:      false,
					Text:           "Create Resource",
					SubmittingText: "Validating...",
					Class:          "mt-4",
				}),
			),
		),
	)
}
