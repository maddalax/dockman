package ui

import (
	"dockside/app/ui/icons"
	"github.com/maddalax/htmgo/framework/h"
)

func SimpleTooltip(children ...h.Ren) *h.Element {
	additionalClasses := ""
	newChildren := make([]h.Ren, 0)

	for _, child := range children {
		skip := false
		if a, ok := child.(*h.AttributeR); ok {
			if a.Name == "class" {
				additionalClasses = a.Value
				skip = true
			}
		}

		if !skip {
			newChildren = append(newChildren, child)
		}
	}

	return h.Div(
		h.Class("has-tooltip h-5 w-5"),
		icons.Question(),
		h.Span(
			h.Class("tooltip rounded shadow-lg p-4 bg-white delay-300", additionalClasses),
			h.Fragment(newChildren...),
		),
	)
}
