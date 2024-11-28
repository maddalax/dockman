package resource

import (
	"dockman/app"
	"dockman/app/ui"
	"dockman/pages/resource/resourceui"
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

	// Get list of all available servers
	servers, err := app.ServerList(locator)

	if err != nil {
		return ui.GenericErrorAlertPartial(ctx, err)
	}

	table := ui.NewTable()
	table2 := ui.NewTable()

	columns := []string{
		"Id",
		"Host Name",
		"IP Address",
		"OS",
		"Last Seen",
		"Status",
		"Actions",
	}

	table.AddColumns(columns)
	table2.AddColumns(columns)

	for _, server := range servers {
		isAssociated := false

		for _, detail := range resource.ServerDetails {
			if detail.ServerId == server.Id {
				isAssociated = true
				break
			}
		}

		t := table
		if !isAssociated {
			t = table2
		}

		t.AddRow()

		runStatus := app.RunStatusNotRunning
		if server.IsAccessible() {
			runStatus = app.RunStatusRunning
		}

		t.AddCellText(server.Id[:8])
		t.AddCellText(server.HostName)
		t.AddCellText(server.IpAddress())
		t.AddCellText(server.Os)
		t.AddCellText(server.LastSeen.Format("2006-01-02 15:04:05"))
		t.AddCell(ui.StatusIndicator(ui.StatusIndicatorProps{
			RunStatus: runStatus,
			TextMap: map[app.RunStatus]string{
				app.RunStatusNotRunning: "Not Accessible",
				app.RunStatusRunning:    "Connected",
			},
		}))

		text := h.Ternary(isAssociated, "Remove from resource", "Associate with resource")
		t.AddCell(
			h.Button(
				h.Text(text),
				h.Class("text-blue-500 hover:text-blue-700"),
				h.GetPartialWithQs(
					ToggleAssociationServerPartial,
					h.NewQs("server_id", server.Id, "resource_id", resource.Id),
					"click",
				),
			),
		)
	}

	// Render the page with two tables
	return h.NewPartial(
		h.Div(
			h.Class("flex flex-col gap-8"),
			h.Id("resource-servers"),
			// Associated Servers Table
			h.Div(
				h.Class("flex flex-col gap-4"),
				h.H2(
					h.Text("Associated Servers"),
				),
				table.Render(),
			),
			// Available Servers Table
			h.Div(
				h.Class("flex flex-col gap-4"),
				h.H2(
					h.Text("Available Servers"),
				),
				table2.Render(),
			),
		),
	)
}
