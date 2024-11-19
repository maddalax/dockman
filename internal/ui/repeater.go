package ui

import (
	"fmt"
	"github.com/maddalax/htmgo/extensions/websocket/ws"
	"github.com/maddalax/htmgo/framework/h"
)

type RepeaterProps struct {
	Item         func(index int) *h.Element
	RemoveButton func(index int, children ...h.Ren) *h.Element
	AddButton    *h.Element
	DefaultItems []*h.Element
	Id           string
	currentIndex int
	OnAdd        func(data ws.HandlerData)
	OnRemove     func(data ws.HandlerData, index int)
}

func (props *RepeaterProps) itemId(index int) string {
	return fmt.Sprintf("%s-repeater-item-%d", props.Id, index)
}

func (props *RepeaterProps) addButtonId() string {
	return fmt.Sprintf("%s-repeater-add-button", props.Id)
}

func repeaterItem(ctx *h.RequestContext, item *h.Element, index int, props *RepeaterProps) *h.Element {
	id := props.itemId(index)
	return h.Div(
		h.Class("flex gap-2 items-center"),
		h.Id(id),
		item,
		props.RemoveButton(
			index,
			ws.OnClick(ctx, func(data ws.HandlerData) {
				if props.OnRemove != nil {
					props.OnRemove(data, index)
				}
				props.currentIndex--
				ws.PushElement(
					data,
					h.Div(
						h.Attribute("hx-swap-oob", fmt.Sprintf("delete:#%s", id)),
						h.Div(),
					),
				)
			}),
		),
	)
}

func Repeater(ctx *h.RequestContext, props RepeaterProps) *h.Element {
	if props.Id == "" {
		props.Id = h.GenId(6)
	}

	props.currentIndex = len(props.DefaultItems)

	return h.Div(
		h.Class("flex flex-col gap-2"),
		h.List(props.DefaultItems, func(item *h.Element, index int) *h.Element {
			return repeaterItem(ctx, item, index, &props)
		}),
		h.Div(
			h.Id(props.addButtonId()),
			h.Class("flex justify-left"),
			props.AddButton,
			ws.OnClick(ctx, func(data ws.HandlerData) {
				if props.OnAdd != nil {
					props.OnAdd(data)
				}
				ws.PushElement(
					data,
					h.Div(
						h.Attribute("hx-swap-oob", "beforebegin:#"+props.addButtonId()),
						repeaterItem(
							ctx, props.Item(props.currentIndex), props.currentIndex, &props,
						),
					),
				)
				props.currentIndex++
			}),
		),
	)
}
