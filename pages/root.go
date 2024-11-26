package pages

import (
	"dockside/__htmgo/assets"
	"fmt"
	"github.com/maddalax/htmgo/extensions/websocket/session"
	"github.com/maddalax/htmgo/framework/h"
)

func WsConnect(ctx *h.RequestContext) *h.AttributeR {
	sessionId := session.GetSessionId(ctx)
	return h.Attribute("ws-connect", fmt.Sprintf("/ws?sessionId=%s", sessionId))
}

func RootPage(ctx *h.RequestContext, children ...h.Ren) *h.Page {
	title := "htmgo template"
	description := "an example of the htmgo template"
	author := "htmgo"
	url := "https://htmgo.dev"

	return h.NewPage(
		h.Html(
			h.Class("h-full bg-white"),
			h.JoinExtensions(
				h.HxExtension(
					h.BaseExtensions(),
				),
				h.HxExtension("ws"),
			),
			h.Head(
				h.Title(
					h.Text(title),
				),
				h.Meta("viewport", "width=device-width, initial-scale=1"),
				h.Link(assets.FaviconIco, "icon"),
				h.Link(assets.AppleTouchIconPng, "apple-touch-icon"),
				h.Meta("title", title),
				h.Meta("charset", "utf-8"),
				h.Meta("author", author),
				h.Meta("description", description),
				h.Meta("og:title", title),
				h.Meta("og:url", url),
				h.Link("canonical", url),
				h.Meta("og:description", description),
				h.Link(assets.MainCss, "stylesheet"),
				h.Script(assets.HtmgoJs),
				h.Tag("script", h.Attribute("src", assets.FloatingUiJs), h.Attribute("type", "module")),
			),
			h.Body(
				h.Class("h-full"),
				WsConnect(ctx),
				h.TriggerChildren(),
				h.Fragment(children...),
			),
		),
	)
}
