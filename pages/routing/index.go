package routing

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"paas/pages"
	"paas/resources"
	"paas/router"
	"paas/slices"
	"paas/ui"
	"paas/ui/icons"
)

func SaveRouteTable(ctx *h.RequestContext) *h.Partial {
	index := 0
	var blocks []router.RouteBlock

	for {
		hostname := ctx.FormValue(fmt.Sprintf("hostname-%d", index))
		path := ctx.FormValue(fmt.Sprintf("path-%d", index))
		resourceId := ctx.FormValue(fmt.Sprintf("resource-%d", index))

		if hostname == "" {
			break
		}

		blocks = append(blocks, router.RouteBlock{
			Hostname:   hostname,
			Path:       path,
			ResourceId: resourceId,
		})

		index++
	}

	// TODO should we automatically apply the blocks here or just save them?
	err := router.ApplyBlocks(ctx.ServiceLocator(), blocks)

	// TODO error handling
	if err != nil {
		return h.NewPartial(h.Div())
	}

	return h.NewPartial(
		h.Div(),
	)
}

func Setup(ctx *h.RequestContext) *h.Page {
	locator := ctx.ServiceLocator()
	resourceNames := resources.GetNames(locator)
	table, err := router.GetRouteTable(locator)

	if err != nil {
		table = []router.RouteBlock{}
	}

	if len(table) == 0 {
		table = []router.RouteBlock{
			{
				Hostname: "",
			},
		}
	}

	fmt.Printf("Route table: %v", table)

	return pages.SidebarPage(
		ctx,
		h.Div(
			h.Class("w-full min-h-screen min-h-[100%] flex flex-col items-center w-full"),
			h.Form(
				h.NoSwap(),
				h.TriggerChildren(),
				h.PostPartial(SaveRouteTable),

				h.Div(
					h.Class("flex justify-between items-center mb-6"),
					h.H2F("Route Table", h.Class("text-xl font-bold")),
					h.Div(
						h.Class("mr-8"),
						ui.PrimaryButton(ui.ButtonProps{
							Text: "Save Changes",
							Type: "submit",
						}),
					),
				),

				ui.Repeater(ctx, ui.RepeaterProps{
					DefaultItems: slices.Map(table, func(rb router.RouteBlock, index int) *h.Element {
						return block(blockProps{
							index:         index,
							path:          rb.Path,
							resourceId:    rb.ResourceId,
							hostname:      rb.Hostname,
							resourceNames: resourceNames,
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
							index:         index,
							resourceNames: resourceNames,
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
	index         int
	hostname      string
	path          string
	resourceId    string
	resourceNames []resources.ResourceName
}

func block(props blockProps) *h.Element {
	return h.Div(
		h.Class("bg-white shadow-md rounded-md p-6 w-full flex gap-6"),
		h.Div(
			h.Class("flex flex-col gap-2"),
			h.Div(
				h.Class("flex gap-1 items-center"),
				// tooltip
				ui.SimpleTooltip(
					h.Class("text-slate-600 text-sm max-w-[300px]"),
					h.Div(
						h.P(h.Text("Hostname Matching"), h.Class("font-bold")),
						h.Class("flex flex-col gap-2"),
						h.P(
							h.Text("Enter the hostname that you want to match, such as 'example.com', `app.example.com`, or `localhost:3000`."),
						),
					),
				),
				// label
				h.Label(h.P(
					h.Text("When "),
					h.Span(h.Text("hostname"), h.Class("font-bold")),
					h.Text(" is"),
				),
				),
			),
			ui.Input(ui.InputProps{
				Name:         fmt.Sprintf("hostname-%d", props.index),
				Placeholder:  "hostname",
				DefaultValue: props.hostname,
				Required:     true,
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
						h.P(h.Text("Path Matching (optional)"), h.Class("font-bold")),
						h.Class("flex flex-col gap-2"),
						h.P(
							h.Text("If you only want to route to this app if its a specific path for that hostname, such as /blog, enter it here."),
						),
						h.P(
							h.Text("Leave it blank to route all paths."),
						),
						h.P(
							h.Text(" The request will match if the path starts with the value you enter, example: /blog will match /blog/my-article."),
						),
					),
				),
				// label
				h.Label(h.P(
					h.Text("When "),
					h.Span(h.Text("path"), h.Class("font-bold")),
					h.Text(" is"),
				),
				),
			),

			h.Div(
				h.Class("flex gap-2 items-center"),
				ui.Input(ui.InputProps{
					Name:         fmt.Sprintf("path-%d", props.index),
					Id:           "path",
					Placeholder:  "(optional) path",
					DefaultValue: props.path,
				}),
			),
		),

		h.Div(
			h.Class("flex flex-col gap-2 min-w-[300px]"),
			h.LabelFor("app-selection", "then route to"),
			ui.Select(ui.SelectProps{
				Id:       fmt.Sprintf("resource-%d", props.index),
				Required: true,
				Value:    props.resourceId,
				Name:     fmt.Sprintf("resource-%d", props.index),
				Items: h.Map(props.resourceNames, func(name resources.ResourceName) ui.Item {
					return ui.Item{
						Value: name.Id,
						Text:  name.Name,
					}
				}),
			}),
		),
	)
}
