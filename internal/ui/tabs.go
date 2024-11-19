package ui

import (
	"github.com/maddalax/htmgo/framework/h"
	"strings"
)

type Link struct {
	Text string
	Href string
}

type LinkTabsProps struct {
	Links []Link
	End   *h.Element
}

func compareLinks(a string, b string) bool {
	a1 := strings.Split(a, "?")
	b1 := strings.Split(b, "?")
	return a1[0] == b1[0]
}

func LinkTabs(ctx *h.RequestContext, props LinkTabsProps) *h.Element {
	currentHref := h.CurrentPath(ctx)
	activeLink := -1
	for i, link := range props.Links {
		if compareLinks(currentHref, link.Href) {
			activeLink = i
			break
		}
	}
	return h.Div(
		h.Div(
			h.Class("sm:hidden"),
			h.Label(
				h.For("Tab"),
				h.Class("sr-only"),
				h.Text("Tab"),
			),
		),
		h.Div(
			h.Div(
				h.Class("border-b border-gray-200 relative"),
				h.Nav(
					h.Class("-mb-px flex gap-6"),
					h.Attribute("aria-label", "LinkTabs"),
					h.List(props.Links, func(link Link, index int) *h.Element {
						return h.A(
							h.Href(link.Href),
							h.Ternary(index == activeLink,
								h.Class("shrink-0 border-b-2 border-brand-500 px-1 pb-4 text-sm font-medium text-brand-600"),
								h.Class("shrink-0 border-b-2 border-transparent px-1 pb-4 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700")),
							h.Text(link.Text),
						)
					}),
				),
				props.End,
			),
		),
	)
}
