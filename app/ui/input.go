package ui

import (
	"dockside/app/ui/icons"
	"github.com/maddalax/htmgo/framework/h"
)

// InputSize represents the available sizes for the input component
type InputSize string

const (
	InputSizeDefault InputSize = "default"
	InputSizeSm      InputSize = "sm"
	InputSizeLg      InputSize = "lg"
)

// InputType represents the available HTML input types
type InputType string

const (
	InputTypeText     InputType = "text"
	InputTypeEmail    InputType = "email"
	InputTypePassword InputType = "password"
	InputTypeSearch   InputType = "search"
	InputTypeFile     InputType = "file"
	InputTypeNumber   InputType = "number"
)

// InputProps defines all possible props for the Input component
type InputProps struct {
	// Core props
	Type        InputType
	Name        string
	Id          string
	Value       string
	Placeholder string
	Required    bool
	Disabled    bool
	ReadOnly    bool
	FullWidth   bool

	// Styling
	Size  InputSize
	Class string

	// Label, description, and help text
	Label       string
	Description string
	HelpText    *h.Element
	Error       string

	// Icons
	LeadingIcon  *h.Element
	TrailingIcon *h.Element

	// HTMX props
	Target  string
	Trigger string
	Get     string
	Post    string

	// Additional props
	Min       string
	Max       string
	Step      string
	Pattern   string
	AutoFocus bool
	Children  []h.Ren
}

// baseInputClasses returns the base classes for the input component
func baseInputClasses() string {
	return "flex rounded-md border border-input bg-background px-3 py-2 text-sm " +
		"ring-offset-background file:border-0 file:bg-transparent " +
		"file:text-sm file:font-medium placeholder:text-muted-foreground " +
		"focus-visible:outline-none focus-visible:ring-2 " +
		"focus-visible:ring-gray-950 dark:focus-visible:ring-gray-300 " +
		"focus-visible:ring-offset-2 disabled:cursor-not-allowed " +
		"disabled:bg-[rgba(0,0,0,0.05)] outline-none focus:ring-0 focus:border-none"
}

// Input creates a new input component with the provided props
func Input(props InputProps) *h.Element {

	if props.Id == "" {
		props.Id = props.Name
	}

	if props.Type == "" {
		props.Type = InputTypeText
	}

	if props.Size == "" {
		props.Size = InputSizeDefault
	}

	// Define size-specific classes
	sizeClasses := map[InputSize]string{
		InputSizeDefault: "h-10",
		InputSizeSm:      "h-8 text-xs",
		InputSizeLg:      "h-12 text-base",
	}

	// Merge all classes
	inputClasses := h.MergeClasses(
		baseInputClasses(),
		sizeClasses[props.Size],
		h.Ternary(props.FullWidth, "w-full", "w-[320px]"),
		// Set default width if not full width
		h.Ternary(props.LeadingIcon != nil, "pl-8", ""),
		h.Ternary(props.TrailingIcon != nil, "pr-8", ""),
		props.Class,
	)

	// Create the input element
	input := h.Input(
		string(props.Type),
		h.Name(props.Name),
		h.Id(props.Id),
		h.Value(props.Value),
		h.Placeholder(props.Placeholder),
		h.Class(inputClasses),
		h.If(
			props.Children != nil,
			h.Children(props.Children...),
		),
		h.If(
			props.Required,
			h.Required(),
		),
		h.If(
			props.Disabled,
			h.Disabled(),
		),
		h.If(
			props.ReadOnly,
			h.ReadOnly(),
		),
		h.If(
			props.AutoFocus,
			h.AutoFocus(),
		),
		h.If(
			props.Min != "",
			h.Min(props.Min),
		),
		h.If(
			props.Max != "",
			h.Max(props.Max),
		),
		h.If(
			props.Step != "",
			h.Step(props.Step),
		),
		h.If(
			props.Pattern != "",
			h.Pattern(props.Pattern),
		),
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
	)

	needsWrapper :=
		props.Label != "" || props.Description != "" || props.HelpText != nil || props.Error != "" ||
			props.LeadingIcon != nil || props.TrailingIcon != nil

	if !needsWrapper {
		return input
	}

	return h.Div(
		h.Class("input-wrapper space-y-1.5"),
		h.If(props.Label != "", FieldLabel(
			props.Label,
			h.For(props.Id),
		)),
		h.If(
			props.Description != "",
			h.P(
				h.Class("text-sm text-muted-foreground"),
				h.Text(props.Description),
			),
		),
		h.Div(
			h.Class(
				h.MergeClasses("relative flex items-center"),
			),
			h.If(
				props.LeadingIcon != nil,
				h.Div(
					h.Class("absolute left-2.5 z-1"),
					props.LeadingIcon,
				),
			),
			input,
			h.If(
				props.TrailingIcon != nil,
				h.Div(
					h.Class("absolute right-2.5 z-1"),
					props.TrailingIcon,
				),
			),
		),
		h.If(
			props.HelpText != nil,
			h.Div(
				h.Class("text-sm text-muted-foreground mt-1"),
				props.HelpText,
			),
		),
		h.If(
			props.Error != "",
			h.P(
				h.Class("text-sm font-medium text-destructive mt-1"),
				h.Text(props.Error),
			),
		),
	)
}

// Helper functions for common input types
func TextInput(props InputProps) *h.Element {
	props.Type = InputTypeText
	return Input(props)
}

func EmailInput(props InputProps) *h.Element {
	props.Type = InputTypeEmail
	return Input(props)
}

func PasswordInput(props InputProps) *h.Element {
	if props.Label == "" {
		props.Label = "Password"
	}
	if props.Name == "" {
		props.Name = "password"
	}
	if props.Placeholder == "" {
		props.Placeholder = "Enter your password"
	}
	props.Type = InputTypePassword
	props.TrailingIcon = icons.EyeIcon()
	return Input(props)
}

func SearchInput(props InputProps) *h.Element {
	props.Type = InputTypeSearch
	props.LeadingIcon = icons.SearchIcon()
	return Input(props)
}

func FileInput(props InputProps) *h.Element {
	props.Type = InputTypeFile
	return Input(props)
}
