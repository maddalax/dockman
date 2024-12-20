package routing

import (
	"dockman/app"
	"dockman/app/ui"
	"dockman/app/ui/icons"
	"dockman/app/util"
	"dockman/pages"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"time"
)

func SaveRouteTable(ctx *h.RequestContext) *h.Partial {
	return util.DelayedPartial(time.Millisecond*800, func() *h.Partial {
		index := 0
		var blocks []app.RouteBlock

		for {
			hostname := ctx.FormValue(fmt.Sprintf("hostname-%d", index))
			path := ctx.FormValue(fmt.Sprintf("path-%d", index))
			resourceId := ctx.FormValue(fmt.Sprintf("resource-%d", index))
			pathMatchModifier := ctx.FormValue(fmt.Sprintf("path-match-modifier-%d", index))

			if hostname == "" {
				break
			}

			blocks = append(blocks, app.RouteBlock{
				Hostname:          hostname,
				Path:              path,
				ResourceId:        resourceId,
				PathMatchModifier: pathMatchModifier,
			})

			index++
		}

		// TODO should we automatically apply the blocks here or just save them?
		err := app.ApplyBlocks(ctx.ServiceLocator(), blocks)

		if err != nil {
			return ui.GenericErrorAlertPartial(ctx, err)
		}

		return ui.SuccessAlertPartial(ctx, "Route Table Saved", "The new rules have been applied.")
	})
}

func Setup(ctx *h.RequestContext) *h.Page {
	locator := ctx.ServiceLocator()
	list, err := app.ResourceList(locator)
	table, err := app.GetRouteTable(locator)

	if err != nil {
		table = []app.RouteBlock{}
	}

	if len(table) == 0 {
		table = []app.RouteBlock{
			{
				Hostname: "",
			},
		}
	}

	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("flex flex-col items-center w-full max-w-5xl mx-auto"),
			h.Form(
				h.Class("w-full"),
				h.NoSwap(),
				h.TriggerChildren(),
				h.PostPartial(SaveRouteTable),
				h.Div(
					h.Class("flex flex-col gap-4 pr-8 mt-6 w-full"),
					ui.AlertPlaceholder(),
					h.Div(
						h.Class("flex justify-between items-center mb-6"),
						h.H2F(
							"Route Table",
							h.Class("text-xl font-bold"),
						),
						h.Div(
							ui.SubmitButton(ui.ButtonProps{
								Text: "Save Changes",
							}),
						),
					),
				),
				ui.Repeater(ctx, ui.RepeaterProps{
					DefaultItems: util.MapSlice(table, func(rb app.RouteBlock, index int) *h.Element {
						return block(blockProps{
							index:             index,
							path:              rb.Path,
							resourceId:        rb.ResourceId,
							pathMatchModifier: rb.PathMatchModifier,
							hostname:          rb.Hostname,
							resources:         list,
						})
					}),
					RemoveButton: func(index int, children ...h.Ren) *h.Element {
						return h.Button(
							h.Type("button"),
							h.Class("w-6 h-6 cursor-pointer"),
							h.Children(children...),
							icons.TrashIcon(),
						)
					},
					Item: func(index int) *h.Element {
						return block(blockProps{
							index:     index,
							resources: list,
						})
					},
					AddButton: h.Div(
						h.Class("mt-1"),
						ui.PrimaryButton(ui.ButtonProps{
							Text:  "+ New Rule",
							Class: "text-sm p-2",
						}),
					),
				}),
			),
		),
	)
}

type blockProps struct {
	index             int
	hostname          string
	path              string
	pathMatchModifier string
	resourceId        string
	resources         []*app.Resource
}

func block(props blockProps) *h.Element {
	return h.Div(
		h.Class("bg-white shadow-md rounded-md p-6 w-full flex flex-col xl:flex-row gap-6 items-center xl:items-start"),
		h.Div(
			h.Class("flex flex-col gap-2 max-w-[350px] w-full"),
			h.Div(
				h.Class("flex gap-1 items-center"),
				// tooltip
				ui.SimpleTooltip(
					h.Class("text-slate-600 text-sm max-w-[300px]"),
					h.Div(
						h.P(
							h.Text("Hostname Matching"),
							h.Class("font-bold"),
						),
						h.Class("flex flex-col gap-2"),
						h.P(
							h.Text("Enter the hostname that you want to match, such as 'example.com', `app.example.com`, or `localhost:3000`."),
						),
					),
				),
				// label
				h.Label(
					h.P(
						h.Text("When "),
						h.Span(
							h.Text("hostname"),
							h.Class("font-bold"),
						),
						h.Text(" is"),
					),
				),
			),
			ui.Input(ui.InputProps{
				Name:        fmt.Sprintf("hostname-%d", props.index),
				Placeholder: "hostname",
				Value:       props.hostname,
				Required:    true,
			}),
		),
		h.Div(
			h.Class("flex flex-col gap-2"),
			h.Div(
				h.Class("flex gap-1 items-center"),
				// tooltip
				ui.SimpleTooltip(
					h.Class("text-slate-600 text-sm max-w-[300px]"),
					h.Div(
						h.P(
							h.Text("Path Matching (optional)"),
							h.Class("font-bold"),
						),
						h.Class("flex flex-col gap-2"),
						h.P(
							h.Text("If you only want to route to this app if its a specific path for that hostname, such as /blog, enter it here."),
						),
						h.P(
							h.Text("Leave it blank to route all paths."),
						),
						h.Div(
							h.P(
								h.Text("Glob Matching"),
								h.Class("font-bold"),
							),
							h.P(
								h.Text(
									"Glob matching is supported. For example, /blog/* will match /blog/my-article, /blog/2021, etc. "),
								h.A(
									h.Text("Learn more"),
									h.Class("text-blue-500 underline"),
									h.Href("https://github.com/gobwas/glob"),
									h.Target("_blank"),
								),
							),
						),
					),
				),
				// label
				h.Label(
					h.P(
						h.Class("flex gap-1 items-center"),
						h.Text("When "),
						h.Span(
							h.Text("path"),
							h.Class("font-bold"),
						),
					),
				),
			),
			h.Div(
				h.Class("flex gap-2 items-center"),
				h.Div(
					h.Class("-mt-1 min-w-[140px]"),
					ui.Select(ui.SelectProps{
						Name:     fmt.Sprintf("path-match-modifier-%d", props.index),
						Required: true,
						Value:    props.pathMatchModifier,
						Items: []ui.Item{
							{
								Value: "equals",
								Text:  "equals",
							},
							{
								Value: "not-equals",
								Text:  "not equals",
							},
							{
								Value: "glob",
								Text:  "glob matches",
							},
							{
								Value: "starts-with",
								Text:  "starts with",
							},
							{
								Value: "not-starts-with",
								Text:  "not starts with",
							},
							{
								Value: "ends-with",
								Text:  "ends with",
							},
							{
								Value: "not-ends-with",
								Text:  "not ends with",
							},
							{
								Value: "contains",
								Text:  "contains",
							},
						},
					}),
				),
				ui.Input(ui.InputProps{
					Name:        fmt.Sprintf("path-%d", props.index),
					Id:          "path",
					Placeholder: "(optional) path",
					Value:       props.path,
				}),
			),
		),
		h.Div(
			h.Class("flex flex-col gap-2 max-w-[350px] w-full"),
			h.LabelFor("app-selection", "then route to"),
			ui.Select(ui.SelectProps{
				Id:       fmt.Sprintf("resource-%d", props.index),
				Required: true,
				Value:    props.resourceId,
				Name:     fmt.Sprintf("resource-%d", props.index),
				Items: h.Map(props.resources, func(name *app.Resource) ui.Item {
					return ui.Item{
						Value: name.Id,
						Text:  name.Name,
					}
				}),
			}),
		),
	)
}
