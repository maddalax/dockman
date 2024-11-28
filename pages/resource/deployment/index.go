package deployment

import (
	"dockside/app"
	"dockside/app/ui"
	"dockside/app/urls"
	"dockside/pages/resource/resourceui"
	"github.com/maddalax/htmgo/framework/h"
)

func Deployment(ctx *h.RequestContext) *h.Page {
	return resourceui.Page(ctx, func(resource *app.Resource) *h.Element {
		return h.Div(
			h.Class("flex flex-col gap-4"),
			h.Div(
				h.Class("flex gap-2 items-center"),
				ui.PrimaryButton(ui.ButtonProps{
					Text: "Start Build",
					Href: urls.ResourceStartDeploymentPath(resource.Id, ""),
				}),
			),
			h.Div(
				h.GetPartialWithQs(
					ListPartial,
					h.NewQs("id", resource.Id),
					"load, every 3s",
				),
			),
		)
	})
}

func ListPartial(ctx *h.RequestContext) *h.Partial {
	deployments, err := app.GetDeployments(ctx.ServiceLocator(), ctx.QueryParam("id"))

	if err != nil {
		deployments = []app.Deployment{}
	}

	table := ui.NewTable()

	table.AddColumns([]string{
		"Build Id",
		"Commit",
		"Ran at",
		"Status",
		"Reason",
		"Actions",
	})

	for _, deployment := range deployments {
		resource, err := app.ResourceGet(ctx.ServiceLocator(), deployment.ResourceId)

		if err != nil {
			continue
		}

		table.AddRow()

		table.AddCellText(deployment.BuildId[:8])

		commitHashUrl := ""

		if deployment.Commit != "" && len(deployment.Commit) > 8 {
			switch bm := resource.BuildMeta.(type) {
			case *app.DockerBuildMeta:
				commitHashUrl = urls.RepoCommitHashUrl(bm.RepositoryUrl, deployment.Commit)
			}
			if commitHashUrl == "" {
				table.AddCellText(deployment.Commit[:8])
			} else {
				table.AddCell(
					h.A(
						h.Href(commitHashUrl),
						h.Text(deployment.Commit[:8]),
						h.Class("text-blue-500 hover:text-blue-700"),
					),
				)
			}
		} else {
			table.AddCell(
				h.Empty(),
			)
		}

		table.WithCellTexts(
			deployment.CreatedAt.Format("Jan 2, 2006 at 3:04 PM"),
			string(deployment.Status),
			deployment.StatusReason,
		)

		table.AddCell(
			h.A(
				h.Href(urls.ResourceDeploymentLogUrl(deployment.ResourceId, deployment.BuildId)),
				h.Text("View Log"),
				h.Class("text-blue-500 hover:text-blue-700"),
			),
		)
	}

	return h.NewPartial(table.Render())
}
