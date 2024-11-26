package resource

import (
	"dockside/app"
	"dockside/app/ui"
	"dockside/pages/resource/resourceui"
	"github.com/maddalax/htmgo/framework/h"
	"slices"
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

	exposedPort, _ := strconv.Atoi(ctx.FormValue("exposed-port"))
	dockerfile := ctx.FormValue("dockerfile")
	deploymentBranch := ctx.FormValue("deployment-branch")
	autoDeploy := ctx.FormValue("auto-deploy") == "on"

	switch bm := resource.BuildMeta.(type) {
	case *app.DockerBuildMeta:
		branches, err := bm.ListRemoteBranches()
		if err == nil && !slices.Contains(branches, deploymentBranch) {
			return ui.ErrorAlertPartial(ctx, h.Pf("Invalid branch"), h.Pf("The deployment branch you specified does not exist in the repository"))
		}
	}

	err = app.ResourcePatch(locator, resource.Id, func(resource *app.Resource) *app.Resource {
		resource.InstancesPerServer = instancesPerServer
		bm := resource.BuildMeta.(*app.DockerBuildMeta)
		bm.DeployOnNewCommit = autoDeploy
		bm.DeploymentBranch = deploymentBranch
		bm.ExposedPort = exposedPort
		bm.Dockerfile = dockerfile
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
				h.Class("flex flex-col gap-5"),
				ui.Input(ui.InputProps{
					Label:    "Resource Name",
					Value:    resource.Name,
					Name:     "name",
					Disabled: true,
				}),
				ui.Input(ui.InputProps{
					Label:    "Resource Type",
					Value:    strconv.Itoa(int(resource.RunType)),
					Disabled: true,
				}),
				ui.Input(ui.InputProps{
					Label:    "Instances Per Server",
					Value:    strconv.Itoa(resource.InstancesPerServer),
					Name:     "instances-per-server",
					HelpText: h.Pf("Number of instances to run on each server, requests will be automatically load balanced between them."),
				}),
				buildMetaFields(resource),
				ui.SubmitButton(ui.ButtonProps{
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
		branches, err := bm.ListRemoteBranches()
		if err != nil {
			branches = []string{}
		}
		return h.Fragment(
			ui.Input(ui.InputProps{
				Label:    "Repository",
				Disabled: true,
				Value:    bm.RepositoryUrl,
				Name:     "repository",
			}),
			h.Div(
				h.Class("flex flex-col gap-1"),
				ui.ComboBox(ui.ComboBoxProps{
					Label:      "Deployment Branch",
					Name:       "deployment-branch",
					Value:      bm.DeploymentBranch,
					ShowSearch: true,
					Items: h.Map(branches, func(item string) h.KeyValue[string] {
						return h.KeyValue[string]{Key: item, Value: item}
					}),
				}),
				ui.Checkbox(ui.CheckboxProps{
					Label:   "Auto Deploy On Push To Branch",
					Checked: bm.DeployOnNewCommit,
					Name:    "auto-deploy",
					Id:      "auto-deploy",
				}),
			),
			//ui.Input(ui.InputProps{
			//	Label:        "Redeploy On Push To Branch",
			//	Value: bm.RedeployOnPushBranch,
			//	Name:         "redeploy-on-push-branch",
			//}),
			ui.Input(ui.InputProps{
				Disabled: true,
				Label:    "Latest Commit",
				Value:    bm.CommitForBuild,
				Name:     "latest-commit",
			}),
			ui.Input(ui.InputProps{
				Label:    "Application Exposed Port",
				Value:    strconv.Itoa(bm.ExposedPort),
				Name:     "exposed-port",
				HelpText: h.Pf("The port your application listens on inside the container, in the case of a docker deployment, its default value is from the EXPOSE directive in the Dockerfile."),
			}),
			ui.Input(ui.InputProps{
				Label: "Dockerfile",
				Value: bm.Dockerfile,
				Name:  "dockerfile",
			}),
		)
	}

	return h.Empty()
}
