package resource

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/pages"
	"paas/ui"
)

func New(ctx *h.RequestContext) *h.Page {
	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("flex flex-col gap-2 w-full items-center"),
			pages.Title("New Resource"),
			ResourceCreateForm(ctx),
		),
	)
}

func ResourceCreateForm(ctx *h.RequestContext) *h.Element {
	return h.Div(
		h.Class("p-6 bg-white w-full"),
		h.Form(
			h.Class("space-y-4 w-full"),
			// Name Input
			ui.Input(ui.InputProps{
				Id:          "name",
				Label:       "Name",
				Name:        "name",
				Required:    true,
				Placeholder: "Enter resource name",
			}),
			// Environment Input
			ui.Input(ui.InputProps{
				Id:          "environment",
				Label:       "Environment",
				Name:        "environment",
				Required:    true,
				Placeholder: "Enter environment",
			}),
			// RunType Dropdown
			h.Div(
				h.Label(
					h.Class("block text-sm font-medium text-gray-700"),
					h.For("run_type"),
					h.Text("Run Type"),
				),
				h.Select(
					h.Class("mt-1 block w-full rounded-md border-gray-300 shadow-sm"),
					h.Id("run_type"),
					h.Name("run_type"),
					h.Required(),
					h.Option(h.Value("Type1"), h.Text("Type1")),
					h.Option(h.Value("Type2"), h.Text("Type2")),
				),
			),
			// BuildMeta Inputs
			ui.Input(ui.InputProps{
				Id:          "dockerfile",
				Label:       "Dockerfile Path",
				Name:        "dockerfile",
				Placeholder: "Enter Dockerfile path",
			}),
			ui.Input(ui.InputProps{
				Id:          "tags",
				Label:       "Tags",
				Name:        "tags",
				Placeholder: "Comma-separated tags",
			}),
			ui.Input(ui.InputProps{
				Id:          "image",
				Label:       "Docker Registry Image",
				Name:        "image",
				Placeholder: "Enter Docker registry image",
			}),
			// Environment Variables TextArea
			h.Div(
				h.Label(
					h.Class("block text-sm font-medium text-gray-700"),
					h.For("env"),
					h.Text("Environment Variables"),
				),
				h.TextArea(
					h.Class("mt-1 block w-full rounded-md border-gray-300 shadow-sm"),
					h.Name("env"),
					h.Id("env"),
					h.Placeholder("Enter environment variables, one per line"),
				),
			),
			// Submit Button using ui.Button
			ui.PrimaryButton(ui.ButtonProps{
				Text:  "Create Resource",
				Type:  "submit",
				Class: "mt-4 w-full",
			}),
		),
	)
}
