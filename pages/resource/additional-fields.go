package resource

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/ui"
)

func AdditionalCreateResourceFields(ctx *h.RequestContext) *h.Partial {
	deploymentType := ctx.QueryParam("deployment_type")
	return h.SwapPartial(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-4"),
			h.Id("additional-create-resource-fields"),
			AdditionalFieldsForDeploymentType(ctx, deploymentType),
			EnvironmentVariables(ctx),
			h.Div(
				h.Id("submit-error"),
			),
			ui.SubmitButton(ui.SubmitButtonProps{
				Text:           "Create Resource",
				SubmittingText: "Validating...",
				Class:          "mt-4 w-full",
			}),
		),
	)
}
