package resource

import (
	"dockside/app"
	"dockside/app/ui"
	"dockside/pages/resource/resourceui"
	"github.com/maddalax/htmgo/framework/h"
	"strconv"
)

func SaveResourceDetails(ctx *h.RequestContext) *h.Partial {
	instancesPerServer, _ := strconv.Atoi(ctx.FormValue("instances-per-server"))
	id := h.GetQueryParam(ctx, "id")

	locator := ctx.ServiceLocator()
	resource, err := app.ResourceGet(locator, id)

	if err != nil {
		return ui.GenericErrorAlertPartial(ctx, err)
	}

	repository := ctx.FormValue("repository")
	redeployOnPushBranch := ctx.FormValue("redeploy-on-push-branch")
	exposedPort, _ := strconv.Atoi(ctx.FormValue("exposed-port"))
	dockerfile := ctx.FormValue("dockerfile")

	err = app.ResourcePatch(locator, resource.Id, func(resource *app.Resource) *app.Resource {
		resource.InstancesPerServer = instancesPerServer
		bm := resource.BuildMeta.(*app.DockerBuildMeta)
		if repository != "" {
			bm.RepositoryUrl = repository
		}
		if redeployOnPushBranch != "" {
			bm.RedeployOnPushBranch = redeployOnPushBranch
		}
		if exposedPort != 0 {
			bm.ExposedPort = exposedPort
		}
		if dockerfile != "" {
			bm.Dockerfile = dockerfile
		}
		resource.BuildMeta = bm
		return resource
	})

	if err != nil {
		return ui.GenericErrorAlertPartial(ctx, err)
	}

	// changed instances per server, start the resource so that the new instances are created or removed
	if resource.InstancesPerServer != instancesPerServer {
		go app.SendResourceStartCommand(locator, resource.Id, app.StartOpts{
			IgnoreIfRunning: true,
			// if we change the instances and existing containers already exist for the new instance indexes, remove them
			RemoveExisting: true,
		})
	}

	return ui.SuccessAlertPartial(ctx, "Resource updated", "Resource details have been updated successfully")
}

func Index(ctx *h.RequestContext) *h.Page {
	return resourceui.Page(ctx, func(resource *app.Resource) *h.Element {
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
				buildMetaFields(resource),
				ui.SubmitButton(ui.SubmitButtonProps{
					Text:           "Save",
					SubmittingText: "Saving...",
					Post:           h.GetPartialPath(SaveResourceDetails),
				}),
			),
		)
	})
}

func buildMetaFields(resource *app.Resource) *h.Element {
	switch bm := resource.BuildMeta.(type) {
	case *app.DockerBuildMeta:
		return h.Fragment(
			ui.Input(ui.InputProps{
				Label:        "Repository",
				Disabled:     true,
				DefaultValue: bm.RepositoryUrl,
				Name:         "repository",
			}),
			ui.Input(ui.InputProps{
				Label:        "Redeploy On Push To Branch",
				DefaultValue: bm.RedeployOnPushBranch,
				Name:         "redeploy-on-push-branch",
			}),
			ui.Input(ui.InputProps{
				Disabled:     true,
				Label:        "Latest Commit",
				DefaultValue: bm.CommitForBuild,
				Name:         "latest-commit",
			}),
			ui.Input(ui.InputProps{
				Label:        "Exposed Port",
				DefaultValue: strconv.Itoa(bm.ExposedPort),
				Name:         "exposed-port",
			}),
			ui.Input(ui.InputProps{
				Label:        "Dockerfile",
				DefaultValue: bm.Dockerfile,
				Name:         "dockerfile",
			}),
		)
	}

	return h.Empty()
}
