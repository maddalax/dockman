package ui

import (
	"dockside/app/ui/icons"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
)

type ComboBoxProps struct {
	Items            []h.KeyValue[string]
	Id               string
	Value            string
	Label            string
	Name             string
	Placeholder      string
	ShowSearch       bool
	UseInput         bool
	Required         bool
	HelpText         *h.Element
	LeadingInputIcon *h.Element
}

func ComboBox(props ComboBoxProps) *h.Element {
	defaultText := "Select an item"

	if props.Id == "" {
		props.Id = h.GenId(6)
	}

	if props.Value != "" {
		selected := h.Find(props.Items, func(item *h.KeyValue[string]) bool {
			return item.Value == props.Value
		})
		if selected != nil {
			defaultText = selected.Key
		}
	}

	if props.UseInput {
		props.ShowSearch = false
	}

	onLoadScriptButton := h.OnLoad(
		// language=JavaScript
		h.EvalJs(`
						  let dropdown = self.nextElementSibling.firstChild;
              let button = self
              let search = dropdown.querySelector('input')
              let valueInput = self.previousElementSibling

 							let stopUpdate = () => window.dockside.floating.stopUpdate(button) 
 							
 							const handleDocClick = (event) => {
								const withinBoundaries = event.composedPath().includes(button) || event.composedPath().includes(dropdown);
								if(!withinBoundaries) {
									 hide()
								}
							}
 							
 							const hide = () => {
                   document.removeEventListener('click', handleDocClick)
                   stopUpdate()
                   setTimeout(() => dropdown.classList.add('hidden'), 25);
 							}
                             
              const show = () => {
                // close dropdown when clicking outside	
   					    document.addEventListener('click', handleDocClick)
                dropdown.classList.remove('hidden');
                window.dockside.floating.updatePosition(button, dropdown);
              }
              
							// dropdown opened
							self.addEventListener("click", () => {
                show()
              })
              
               // close dropdown when clicking an option
							dropdown.querySelectorAll('li').forEach(li => {
                 li.addEventListener('click', event => {
										event.stopPropagation()
										const target = event.target.tagName === 'LI' ? event.target : event.target.parentElement
										if(target.innerText) {
												button.querySelector('[data-value="label"]').innerText = target.innerText;   
												valueInput.value = target.getAttribute('data-value')
												hide()  
										}
									});
							})
							
							// handle search
							if(search) {
								 search.addEventListener('input', event => {
											const value = event.target.value.toLowerCase();
											dropdown.querySelectorAll('li').forEach(li => {
													const text = li.innerText.toLowerCase();
													if(text.includes(value)) {
															li.classList.remove('hidden');
													} else {
															li.classList.add('hidden');
													}
											})
								})                                  
							}
            `),
	)

	onLoadScriptInput := h.OnLoad(
		// language=JavaScript
		h.EvalJs(fmt.Sprintf(`					
						  let dropdown = document.getElementById('%s-combobox-options');
              let input = self

 							let stopUpdate = () => window.dockside.floating.stopUpdate(input) 
 							
 							const handleDocClick = (event) => {
								const withinBoundaries = event.composedPath().includes(input) || event.composedPath().includes(dropdown);
								if(!withinBoundaries) {
									 hide()
								}
							}
 							
 							const hide = () => {
							   document.removeEventListener('click', handleDocClick)
							   stopUpdate()
							   setTimeout(() => dropdown.classList.add('hidden'), 25);
 							}
                             
              const show = () => {
                // close dropdown when clicking outside	
   					    document.addEventListener('click', handleDocClick)
                dropdown.classList.remove('hidden');
                window.dockside.floating.updatePosition(input, dropdown);
              }
              
							// dropdown opened
							self.addEventListener("focus", () => {
                show()
              })
              
               // close dropdown when clicking an option
							dropdown.querySelectorAll('li').forEach(li => {
                 li.addEventListener('click', event => {
										event.stopPropagation()
										const target = event.target.tagName === 'LI' ? event.target : event.target.parentElement
										if(target.innerText) {
												input.innerText = target.innerText;   
												input.value = target.getAttribute('data-value')
												hide()  
										}
									});
							})
							
							// handle search
							input.addEventListener('input', event => {
											const value = event.target.value.toLowerCase();
                                            let empty = true;
											dropdown.querySelectorAll('li').forEach(li => {
													const text = li.innerText.toLowerCase();
													if(text.includes(value)) {
															li.classList.remove('hidden');
                                                            empty = false;
													} else {
															li.classList.add('hidden');
													}
											})
											if(empty) {
                                                dropdown.classList.add('hidden');
											} else {
                                                dropdown.classList.remove('hidden');
											}
								})   
            `, props.Id)),
	)

	dropdown := h.Div(
		h.Class("relative max-w-[320px]"),
		h.Div(
			h.Id(fmt.Sprintf("%s-combobox-options", props.Id)),
			h.Role("listbox"),
			h.Class("hidden absolute z-40 rounded-md border bg-popover text-popover-foreground shadow-lg outline-none animate-in fade-in-0 zoom-in-95 max-w-[320px]"),
			h.If(props.ShowSearch, SearchInput(InputProps{
				AutoFocus:   true,
				Placeholder: "Search...",
				FullWidth:   true,
			})),
			h.Ul(
				h.Class("max-h-60 overflow-auto p-1 w-[320px]"),
				h.List(props.Items, func(item h.KeyValue[string], index int) *h.Element {
					return h.Li(
						h.Role("option"),
						h.Attribute("aria-selected", "false"),
						h.Class("relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-4 pr-2 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 hover:bg-accent hover:text-accent-foreground"),
						h.Span(
							h.Class("flex justify-between w-full"),
							h.Text(item.Key),
						),
						h.Attribute("data-value", item.Value),
					)
				}),
			),
		),
	)

	comboboxInput := h.Div(
		h.Class("w-full max-w-[320px]"),
		Input(InputProps{
			Name:        props.Name,
			Type:        InputTypeSearch,
			HelpText:    props.HelpText,
			LeadingIcon: props.LeadingInputIcon,
			Required:    props.Required,
			Placeholder: props.Placeholder,
			Value:       props.Value,
			Children: []h.Ren{
				onLoadScriptInput,
			},
		}),
		dropdown,
	)

	comboboxButton := h.Div(
		h.Class("w-full max-w-[320px]"),
		h.Input(
			"text",
			h.Name(props.Name),
			h.Value(props.Value),
			h.Class("hidden"),
			h.If(
				props.Required,
				h.Required(),
			),
		),
		h.Button(
			h.Type("button"),
			h.Role("combobox"),
			h.Attribute("aria-controls", "combobox-options"),
			h.Attribute("aria-expanded", "true"),
			h.Attribute("aria-autocomplete", "list"),
			h.Class("flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"),
			h.Span(
				h.Class("truncate"),
				h.Attribute("data-value", "label"),
				h.Text(defaultText),
			),
			icons.ChevronDown(),
			onLoadScriptButton,
		),
		dropdown,
	)

	comp := h.Ternary(props.UseInput, comboboxInput, comboboxButton)

	if props.Label == "" {
		return comp
	}

	return h.Div(
		h.Class("flex flex-col space-y-2"),
		h.Label(
			h.Class("text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"),
			h.For(props.Name),
			h.Text(props.Label),
		),
		comp,
	)
}
