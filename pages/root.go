package pages

import (
	"fmt"
	"github.com/maddalax/htmgo/extensions/websocket/session"
	"github.com/maddalax/htmgo/framework/h"
	"paas/__htmgo/assets"
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
			),
			h.Body(
				WsConnect(ctx),
				h.Div(
					h.Class("flex flex-col gap-2 bg-white h-full"),
					h.Fragment(children...),
				),
			),
		),
	)
}
