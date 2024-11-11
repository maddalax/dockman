package resource

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"paas/pages"
	"paas/resources"
	"paas/ui"
)

func New(ctx *h.RequestContext) *h.Page {
	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-2 w-full items-center"),
			pages.Title("New Resource"),
			CreateForm(ctx),
		),
	)
}

func CreateForm(ctx *h.RequestContext) *h.Element {
	return h.Div(
		h.Class("p-6 bg-white w-full"),
		h.Form(
			h.PostPartial(SubmitHandler),
			h.Class("space-y-4 w-full"),

			ui.Input(ui.InputProps{
				Id:          "name",
				Label:       "Name",
				Name:        "name",
				Required:    true,
				Placeholder: "Enter resource name",
			}),

			ui.Input(ui.InputProps{
				Id:          "environment",
				Label:       "Environment",
				Name:        "environment",
				Required:    true,
				Placeholder: "Enter environment",
			}),

			DeploymentChoiceSelector(),

			AdditionalFieldsForDeploymentType(ctx, ""),

			//// BuildMeta Inputs
			//ui.Input(ui.InputProps{
			//	Id:          "dockerfile",
			//	Label:       "Dockerfile Path",
			//	Name:        "dockerfile",
			//	Placeholder: "Enter Dockerfile path",
			//}),
			//ui.Input(ui.InputProps{
			//	Id:          "tags",
			//	Label:       "Tags",
			//	Name:        "tags",
			//	Placeholder: "Comma-separated tags",
			//}),
			//ui.Input(ui.InputProps{
			//	Id:          "image",
			//	Label:       "Docker Registry Image",
			//	Name:        "image",
			//	Placeholder: "Enter Docker registry image",
			//}),
			//// Environment Variables TextArea
			//h.Div(
			//	h.Label(
			//		h.Class("block text-sm font-medium text-gray-700"),
			//		h.For("env"),
			//		h.Text("Environment Variables"),
			//	),
			//	h.TextArea(
			//		h.Class("mt-1 block w-full rounded-md border-gray-300 shadow-sm"),
			//		h.Name("env"),
			//		h.Id("env"),
			//		h.Placeholder("Enter environment variables, one per line"),
			//	),
			//),
			// Submit Button using ui.Button
			//ui.PrimaryButton(ui.ButtonProps{
			//	Text:  "Create Resource",
			//	Type:  "submit",
			//	Class: "mt-4 w-full",
			//}),
		),
	)
}

func SubmitHandler(ctx *h.RequestContext) *h.Partial {
	//data := ctx.FormValue("data")
	ctx.Request.ParseForm()
	values := ctx.Request.Form
	fmt.Println(values)

	env := make(map[string]string)
	for i := 0; i < 100; i++ {
		key := values.Get(fmt.Sprintf("env-key-%d", i))
		value := values.Get(fmt.Sprintf("env-value-%d", i))
		if key == "" || value == "" {
			break
		}
		env[key] = value
	}

	runType := resources.Unknown

	if values.Get("deployment-type") == "dockerfile" {
		runType = resources.RunTypeDockerBuild
	}

	err := resources.Create(ctx.ServiceLocator(), resources.CreateOptions{
		Name:        values.Get("name"),
		Environment: values.Get("environment"),
		RunType:     runType,
		BuildMeta: resources.DockerBuildMeta{
			RepositoryUrl:     values.Get("git-repository"),
			Dockerfile:        values.Get("dockerfile"),
			GithubAccessToken: values.Get("github-access-token"),
			Tags:              []string{},
		},
		Env: env,
	})

	if err != nil {
		return h.NewPartial(h.Div(
			h.Class("text-red-500"),
			h.Text(err.Error()),
		))
	}

	return h.NewPartial(h.Div())
}
