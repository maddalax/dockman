package resource

import (
	"dockside/app"
	"dockside/app/ui"
	"dockside/pages/resource/resourceui"
	"errors"
	"github.com/maddalax/htmgo/framework/h"
)

func ToggleAssociationServerPartial(ctx *h.RequestContext) *h.Partial {
	serverId := ctx.QueryParam("server_id")
	resourceId := ctx.QueryParam("resource_id")
	locator := ctx.ServiceLocator()
	resource, err := app.ResourceGet(locator, resourceId)

	if serverId == "" || resourceId == "" || err != nil {
		return ui.GenericErrorAlertPartial(ctx, errors.New("invalid server or resource"))
	}

	isAssociated := false

	for _, detail := range resource.ServerDetails {
		if detail.ServerId == serverId {
			isAssociated = true
			break
		}
	}

	if isAssociated {
		err = app.DetachServerFromResource(locator, serverId, resourceId)
	} else {
		err = app.AttachServerToResource(locator, serverId, resourceId)
	}

	if err != nil {
		return ui.GenericErrorAlertPartial(ctx, err)
	}

	return h.SwapManyPartial(ctx,
		ui.SuccessAlert(
			h.Pf(""),
			h.Pf("Server associated successfully"),
		),
		ServerListPartial(ctx).Root,
	)
}

func ServerPage(ctx *h.RequestContext) *h.Page {
	return resourceui.Page(ctx, func(resource *app.Resource) *h.Element {
		return h.Div(
			h.GetPartialWithQs(
				ServerListPartial,
				h.NewQs("resource_id", resource.Id),
				"load, every 3s",
			),
		)
	})
}

func ServerListPartial(ctx *h.RequestContext) *h.Partial {
	locator := ctx.ServiceLocator()
	resourceId := ctx.QueryParam("resource_id")
	resource, err := app.ResourceGet(locator, resourceId)

	if err != nil {
		return ui.GenericErrorAlertPartial(ctx, err)
	}

	// Get detailed list of associated servers
	associatedDetails := h.Map(resource.ServerDetails, func(details app.ResourceServer) app.ResourceServerWithDetails {
		server, err := app.ServerGet(locator, details.ServerId)
		if err != nil {
			return app.ResourceServerWithDetails{}
		}
		return app.ResourceServerWithDetails{
			ResourceServer: &details,
			Details:        server,
		}
	})

	// Get list of all available servers
	allServers, err := app.ServerList(locator)

	if err != nil {
		return ui.GenericErrorAlertPartial(ctx, err)
	}

	// Identify unassociated servers
	associatedServerIds := make(map[string]bool)

	for _, server := range associatedDetails {
		associatedServerIds[server.Details.Id] = true
	}

	unassociatedServers := h.Filter(allServers, func(server *app.Server) bool {
		return !associatedServerIds[server.Id]
	})

	// Render the page with two tables
	return h.NewPartial(
		h.Div(
			h.Class("flex flex-col gap-8"),
			h.Id("resource-servers"),
			// Associated Servers Table
			h.Div(
				h.H2(
					h.Text("Associated Servers"),
					h.Class("text-xl font-bold mb-4"),
				),
				renderServerTable(
					associatedDetails,
					func(server app.ResourceServerWithDetails, index int) *h.Element {
						return serverBlockRow(server.Details, resource, true)
					},
				),
			),
			// Unassociated Servers Table
			h.Div(
				h.H2(
					h.Text("Available Servers"),
					h.Class("text-xl font-bold mb-4"),
				),
				renderServerTable(
					unassociatedServers,
					func(server *app.Server, index int) *h.Element {
						return serverBlockRow(server, resource, false)
					},
				),
			),
		),
	)
}

// Helper to render a table
func renderServerTable[T any](servers []T, rowRenderer func(server T, index int) *h.Element) *h.Element {
	return h.Table(
		h.Class("w-full table-auto border-collapse border border-gray-200"),
		h.THead(
			h.Tr(
				h.Th(
					h.Text("Id"),
					h.Class("border border-gray-200 px-4 py-2"),
				),
				h.Th(
					h.Text("Host Name"),
					h.Class("border border-gray-200 px-4 py-2"),
				),
				h.Th(
					h.Text("IP Address"),
					h.Class("border border-gray-200 px-4 py-2"),
				),
				h.Th(
					h.Text("OS"),
					h.Class("border border-gray-200 px-4 py-2"),
				),
				h.Th(
					h.Text("Last Seen"),
					h.Class("border border-gray-200 px-4 py-2"),
				),
				h.Th(
					h.Text("Status"),
					h.Class("border border-gray-200 px-4 py-2"),
				),
				h.Th(
					h.Text("Actions"),
					h.Class("border border-gray-200 px-4 py-2"),
				),
			),
		),
		h.TBody(
			h.List(servers, rowRenderer),
		),
	)
}

// Helper to render a server block as a row
func serverBlockRow(server *app.Server, resource *app.Resource, isAssociated bool) *h.Element {
	runStatus := app.RunStatusNotRunning
	if server.IsAccessible() {
		runStatus = app.RunStatusRunning
	}

	return h.Tr(
		h.Td(
			h.Text(server.Id),
			h.Class("border border-gray-200 px-4 py-2"),
		),
		h.Td(
			h.Text(server.HostName),
			h.Class("border border-gray-200 px-4 py-2"),
		),
		h.Td(
			h.Text(server.IpAddress()),
			h.Class("border border-gray-200 px-4 py-2"),
		),
		h.Td(
			h.Text(server.Os),
			h.Class("border border-gray-200 px-4 py-2"),
		),
		h.Td(
			h.Text(server.LastSeen.Format("2006-01-02 15:04:05")),
			h.Class("border border-gray-200 px-4 py-2"),
		),
		h.Td(
			ui.StatusIndicator(ui.StatusIndicatorProps{
				RunStatus: runStatus,
				TextMap: map[app.RunStatus]string{
					app.RunStatusNotRunning: "Not Accessible",
					app.RunStatusRunning:    "Connected",
				},
			}),
			h.Class("border border-gray-200 px-4 py-2"),
		),
		h.Td(
			h.IfElse(isAssociated,
				ui.PrimaryButton(ui.ButtonProps{
					Text: "Remove from resource",
					Post: h.GetPartialPathWithQs(
						ToggleAssociationServerPartial,
						h.NewQs("server_id", server.Id, "resource_id", resource.Id),
					),
				}),
				ui.PrimaryButton(ui.ButtonProps{
					Text: "Associate with resource",
					Post: h.GetPartialPathWithQs(
						ToggleAssociationServerPartial,
						h.NewQs("server_id", server.Id, "resource_id", resource.Id),
					),
				})),
		),
	)
}
