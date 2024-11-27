package ui

import "github.com/maddalax/htmgo/framework/h"

func FieldLabel(label string, children ...h.Ren) *h.Element {
	classes := "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
	return h.Label(
		h.Class(classes),
		h.Text(label),
		h.Children(children...),
	)
}
