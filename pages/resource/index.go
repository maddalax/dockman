package resource

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/internal"
	"paas/internal/ui"
	"paas/pages/resource/resourceui"
	"strconv"
)

func SaveResourceDetails(ctx *h.RequestContext) *h.Partial {
	instancesPerServer, _ := strconv.Atoi(ctx.FormValue("instances-per-server"))
	id := h.GetQueryParam(ctx, "id")

	locator := ctx.ServiceLocator()
	resource, err := internal.ResourceGet(locator, id)

	if err != nil {
		return ui.GenericErrorAlertPartial(ctx, err)
	}

	err = internal.ResourcePatch(locator, resource.Id, func(resource *internal.Resource) *internal.Resource {
		resource.InstancesPerServer = instancesPerServer
		return resource
	})

	if err != nil {
		return ui.GenericErrorAlertPartial(ctx, err)
	}

	// changed instances per server, start the resource
	if resource.InstancesPerServer != instancesPerServer {
		go internal.Start(locator, resource.Id, internal.StartOpts{
			IgnoreIfRunning: true,
			// if we change the instances and existing containers already exist for the new instance indexes, remove them
			RemoveExisting: true,
		})
	}

	return ui.SuccessAlertPartial(ctx, "Resource updated", "Resource details have been updated successfully")
}

func Index(ctx *h.RequestContext) *h.Page {
	return resourceui.Page(ctx, func(resource *internal.Resource) *h.Element {
		return h.Div(
			h.Class("flex flex-col gap-4"),
			ui.AlertPlaceholder(),
			h.Form(
				h.NoSwap(),
				h.Class("flex flex-col gap-2"),
				ui.Input(ui.InputProps{
					Label:        "Resource Name",
					DefaultValue: resource.Name,
					Name:         "name",
					Disabled:     true,
				}),
				ui.Input(ui.InputProps{
					Label:        "Resource Type",
					DefaultValue: strconv.Itoa(int(resource.RunType)),
					Disabled:     true,
				}),
				ui.Input(ui.InputProps{
					Label:        "Instances Per Server",
					DefaultValue: strconv.Itoa(resource.InstancesPerServer),
					Name:         "instances-per-server",
					HelpText:     h.Pf("Number of instances to run on each server, requests will be automatically load balanced between them."),
				}),
				ui.SubmitButton(ui.SubmitButtonProps{
					Text:           "Save",
					SubmittingText: "Saving...",
					Post:           h.GetPartialPath(SaveResourceDetails),
				}),
			),
		)
	})
}
