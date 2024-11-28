package ui

import (
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
)

type ButtonSize string
type ButtonVariant string

const (
	ButtonSizeXs ButtonSize = "xs"
	ButtonSizeSm ButtonSize = "sm"
	ButtonSizeMd ButtonSize = "md"
	ButtonSizeLg ButtonSize = "lg"
	ButtonSizeXl ButtonSize = "xl"
)

const (
	ButtonVariantDefault     ButtonVariant = "default"
	ButtonVariantPrimary     ButtonVariant = "primary"
	ButtonVariantSecondary   ButtonVariant = "secondary"
	ButtonVariantDestructive ButtonVariant = "destructive"
	ButtonVariantGhost       ButtonVariant = "ghost"
	ButtonVariantLink        ButtonVariant = "link"
)

type ButtonProps struct {
	// Core props
	Text      string
	Disabled  bool
	FullWidth bool

	// Styling
	Size    ButtonSize
	Variant ButtonVariant
	Class   string

	// Icons
	LeftIcon  *h.Element
	RightIcon *h.Element

	// HTMX and interaction props
	Target   string
	Type     string
	Trigger  string
	Get      string
	Post     string
	Href     string
	Children []h.Ren

	ShowLoader bool
}

func Button(props ButtonProps) *h.Element {
	baseClasses := "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium " +
		"ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 " +
		"focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 " +
		"[&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0"

	sizeClasses := map[ButtonSize]string{
		ButtonSizeXs: "h-8 px-3 text-xs",
		ButtonSizeSm: "h-9 px-3",
		ButtonSizeMd: "h-10 px-4 py-2",
		ButtonSizeLg: "h-11 px-8",
		ButtonSizeXl: "h-12 px-8",
	}

	variantClasses := map[ButtonVariant]string{
		ButtonVariantDefault:     "bg-primary text-primary-foreground hover:bg-primary/90",
		ButtonVariantPrimary:     "bg-primary text-primary-foreground hover:bg-primary/90",
		ButtonVariantSecondary:   "bg-secondary text-secondary-foreground hover:bg-secondary/80",
		ButtonVariantDestructive: "bg-destructive text-destructive-foreground hover:bg-destructive/90",
		ButtonVariantGhost:       "hover:bg-accent hover:text-accent-foreground",
		ButtonVariantLink:        "text-primary underline-offset-4 hover:underline",
	}

	if props.Size == "" {
		props.Size = ButtonSizeMd
	}

	if props.Variant == "" {
		props.Variant = ButtonVariantDefault
	}

	if props.ShowLoader {
		appendLoaderProps(&props)
	}

	classes := h.MergeClasses(
		baseClasses,
		sizeClasses[props.Size],
		variantClasses[props.Variant],
		h.Ternary(props.FullWidth, "w-full", "w-auto"),
		props.Class,
	)

	tag := h.Ternary(props.Href != "", "a", "button")

	children := make([]h.Ren, 0)

	if props.LeftIcon != nil {
		children = append(children, props.LeftIcon)
	}

	if props.Text != "" {
		children = append(
			children,
			h.Text(props.Text),
		)
	}

	if props.RightIcon != nil {
		children = append(children, props.RightIcon)
	}

	props.Children = append(props.Children, children...)

	return h.Tag(
		tag,
		h.Class(classes),
		h.If(
			props.Target != "",
			h.HxTarget(props.Target),
		),
		h.If(
			props.Trigger != "",
			h.HxTriggerString(props.Trigger),
		),
		h.If(
			props.Get != "",
			h.Get(props.Get),
		),
		h.If(
			props.Post != "",
			h.Post(props.Post),
		),
		h.If(
			props.Href != "",
			h.Href(props.Href),
		),
		h.IfElse(
			props.Type != "",
			h.Type(props.Type),
			h.Type("button"),
		),
		h.If(
			props.Disabled,
			h.Disabled(),
		),
		h.Children(props.Children...),
	)
}

// Helper functions for variant buttons
func DefaultButton(props ButtonProps) *h.Element {
	props.Variant = ButtonVariantDefault
	return Button(props)
}

func PrimaryButton(props ButtonProps) *h.Element {
	props.Variant = ButtonVariantPrimary
	return Button(props)
}

func SecondaryButton(props ButtonProps) *h.Element {
	props.Variant = ButtonVariantSecondary
	return Button(props)
}

func DestructiveButton(props ButtonProps) *h.Element {
	props.Variant = ButtonVariantDestructive
	return Button(props)
}

func GhostButton(props ButtonProps) *h.Element {
	props.Variant = ButtonVariantGhost
	return Button(props)
}

func LinkButton(props ButtonProps) *h.Element {
	props.Variant = ButtonVariantLink
	return Button(props)
}

func DangerButton(props ButtonProps) *h.Element {
	props.Variant = ButtonVariantDestructive
	return Button(props)
}

func SubmitButton(props ButtonProps) *h.Element {
	props.Type = "submit"
	props.ShowLoader = true
	return Button(props)
}

func appendLoaderProps(props *ButtonProps) {
	props.LeftIcon = spinner()

	if props.Children == nil {
		props.Children = make([]h.Ren, 0)
	}

	props.Children = append(
		props.Children,
		h.OnLoad(
			// language=JavaScript
			h.EvalJs(`
		  let form = self.closest('form');
      		let startLoader = new Function(self.dataset.startLoader);
			if(form) {
			   form.addEventListener('submit', function() {
				 startLoader.call(self);
			   })
      		} else {
			 // if the button is not in a form, we need to manually trigger the click event
			   self.addEventListener('click', function() {
			   startLoader.call(self);
			 })
      }
     `),
		),
	)

	props.Children = append(
		props.Children,
		h.OnEvent(
			"data-start-loader",
			js.RunAfterTimeout(
				100,
				js.AddAttribute("disabled", ""),
				js.RemoveClassOnChildren(".spinner", "hidden"),
			),
		),
	)

	props.Children = append(
		props.Children,
		h.HxAfterRequest(
			// delay so the loading spinner doesn't flash too quickly
			// and we give some feedback to the user
			js.RunAfterTimeout(
				100,
				js.RemoveAttribute("disabled"),
				js.SetClassOnChildren(".spinner", "hidden"),
			),
		),
	)
}

func getSizeIconClass(size ButtonSize) string {
	switch size {
	case ButtonSizeXs:
		return "h-3 w-3"
	case ButtonSizeSm:
		return "h-4 w-4"
	case ButtonSizeMd:
		return "h-5 w-5"
	case ButtonSizeLg:
		return "h-5 w-5"
	case ButtonSizeXl:
		return "h-6 w-6"
	default:
		return "h-5 w-5"
	}
}

func spinner(children ...h.Ren) *h.Element {
	return h.Div(
		h.Children(children...),
		h.Class("hidden spinner spinner-border animate-spin inline-block w-4 h-4 border-2 rounded-full border-slate-200 border-t-transparent"),
		h.Attribute("role", "status"),
	)
}
