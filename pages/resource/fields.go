package resource

import (
	"fmt"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/hx"
	"github.com/maddalax/htmgo/framework/js"
	"paas/ui"
	"paas/ui/icons"
)

func EnvironmentVariables(ctx *h.RequestContext) *h.Element {

	var item = func(index int) *h.Element {
		return h.Div(
			h.Class("flex gap-2"),
			ui.Input(ui.InputProps{
				Id:          "key",
				Label:       h.Ternary(index == 0, "Name", ""),
				Name:        fmt.Sprintf("env-key-%d", index),
				Placeholder: "ENV",
			}),
			ui.Input(ui.InputProps{
				Id:          "value",
				Label:       h.Ternary(index == 0, "Value", ""),
				Name:        fmt.Sprintf("env-value-%d", index),
				Placeholder: "production",
			}),
		)
	}

	return h.Div(
		h.Class("flex flex-col gap-2"),
		h.Label(h.Text("Environment Variables")),
		ui.Repeater(ctx, ui.RepeaterProps{
			DefaultItems: []*h.Element{
				item(0),
			},
			Id: "environment-variables",
			OnAdd: func(data ws.HandlerData) {
				//ws.BroadcastServerSideEvent("increment", map[string]any{})
			},
			OnRemove: func(data ws.HandlerData, index int) {
				//ws.BroadcastServerSideEvent("decrement", map[string]any{})
			},
			AddButton: h.Button(
				h.Type("button"),
				h.Text("+ New Environment Variable"),
			),
			RemoveButton: func(index int, children ...h.Ren) *h.Element {
				if index == 0 {
					return nil
				}
				return h.Button(
					h.Type("button"),
					h.Class("w-6 h-6 cursor-pointer"),
					icons.TrashIcon(),
					h.Children(children...),
				)
			},
			Item: func(index int) *h.Element {
				return item(index)
			},
		}),
	)
}

func AdditionalFieldsForDeploymentType(ctx *h.RequestContext, deploymentType string) *h.Element {
	switch deploymentType {
	case "dockerfile":
		return h.Div(
			h.Class("flex flex-col gap-4"),
			ui.Input(ui.InputProps{
				Id:          "git-repository",
				Label:       "Git Repository Url",
				Name:        "git-repository",
				Placeholder: "https://github.com/maddalax/paas",
				Children: []h.Ren{
					h.OnEvent(hx.KeyUpEvent, js.EvalJs(
						// language=JavaScript
						`
           let next = document.getElementById("git-access-token-input");
           let isGithub = self.value.toLowerCase().includes("github.com/");
           isGithub ? next.classList.remove("hidden") : next.classList.add("hidden");
					`)),
				},
			}),
			h.Div(
				h.Id("git-access-token-input"),
				h.Class("hidden"),
				ui.Input(ui.InputProps{
					Id:          "git-access-token",
					Label:       "Github Repository Access Token (optional)",
					Name:        "github-access-token",
					Placeholder: "",
					HelpText: h.Fragment(
						h.P(
							h.Text("If this is a private repository, provide a git personal access token so the repository can be cloned. "),
							h.A(
								h.Class("text-brand-500 underline"),
								h.Text("More Info"),
								h.Target("_blank"),
								h.Href("https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token"),
							),
							h.Text("."),
						),
						h.P(
							h.Text("Ensure the token has the 'Contents' repository permission."),
						),
					),
				}),
			),
			ui.Input(ui.InputProps{
				Id:          "dockerfile",
				Label:       "Dockerfile Path",
				Name:        "dockerfile",
				Placeholder: "./app/Dockerfile",
				Required:    true,
				HelpText:    h.Pf("The path to the Dockerfile relative to the root of the repository"),
			}),
		)

	case "docker-registry":
		return h.Div()
	}

	return h.Div(
		h.Id("additional-create-resource-fields"),
	)
}
