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
			ui.PrimaryButton(ui.ButtonProps{
				Text:  "Create Resource",
				Type:  "submit",
				Class: "mt-4 w-full",
			}),
		),
	)
}
