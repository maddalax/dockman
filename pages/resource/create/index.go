package create

import (
	"dockside/app"
	"dockside/app/ui"
	"dockside/app/urls"
	"dockside/pages"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
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
			h.TriggerChildren(),
			h.NoSwap(),
			h.PostPartial(SubmitHandler),
			h.Class("space-y-4 w-full"),
			ui.Input(ui.InputProps{
				Id:          "name",
				Label:       "Name",
				Name:        "name",
				Required:    true,
				Placeholder: "Enter resource name",
			}),
			EnvironmentInput(ctx),
			DeploymentChoiceSelector(),
			AdditionalFieldsForDeploymentType(ctx, ""),
		),
	)
}

func SubmitHandler(ctx *h.RequestContext) *h.Partial {
	ctx.Request.ParseForm()
	values := ctx.Request.Form

	env := make(map[string]string)

	index := 0
	for {
		key := values.Get(fmt.Sprintf("env-key-%d", index))
		value := values.Get(fmt.Sprintf("env-value-%d", index))
		if key == "" || value == "" {
			break
		}
		env[key] = value
		index++
	}

	runType := app.RunTypeUnknown

	if values.Get("deployment-type") == "dockerfile" {
		runType = app.RunTypeDockerBuild
	}

	var createBuildMeta = func() app.BuildMeta {
		if runType == app.RunTypeDockerBuild {
			return &app.DockerBuildMeta{
				RepositoryUrl:     values.Get("git-repository"),
				Dockerfile:        values.Get("dockerfile"),
				GithubAccessToken: values.Get("github-access-token"),
				Tags:              []string{},
			}
		}
		return &app.EmptyBuildMeta{}
	}

	id, err := app.ResourceCreate(ctx.ServiceLocator(), app.ResourceCreateOptions{
		Name:        values.Get("name"),
		Environment: values.Get("environment"),
		RunType:     runType,
		BuildMeta:   createBuildMeta(),
		Env:         env,
	})

	if err != nil {
		return h.SwapPartial(
			ctx,
			h.Div(
				h.Id("submit-error"),
				ui.ErrorAlert(
					h.Pf("Unable to create resource"),
					h.Pf(err.Error()),
				),
			),
		)
	}

	return h.RedirectPartial(
		urls.ResourceUrl(id),
	)
}
