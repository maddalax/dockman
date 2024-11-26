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
	Size      InputSize
	Class     string
	WrapClass string

	// Label, description, and help text
	Label            string
	LabelClass       string
	Description      string
	DescriptionClass string
	HelpText         *h.Element
	Error            string
	ErrorClass       string

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
		"disabled:bg-[rgba(0,0,0,0.05)] outline-none"
}

// Input creates a new input component with the provided props
func Input(props InputProps) *h.Element {
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
		h.Ternary(props.FullWidth, "w-full", "w-[320px]"), // Set default width if not full width
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
		h.If(props.Children != nil, h.Children(props.Children...)),
		h.If(props.Required, h.Required()),
		h.If(props.Disabled, h.Disabled()),
		h.If(props.ReadOnly, h.ReadOnly()),
		h.If(props.AutoFocus, h.AutoFocus()),
		h.If(props.Min != "", h.Min(props.Min)),
		h.If(props.Max != "", h.Max(props.Max)),
		h.If(props.Step != "", h.Step(props.Step)),
		h.If(props.Pattern != "", h.Pattern(props.Pattern)),
		h.If(props.Target != "", h.HxTarget(props.Target)),
		h.If(props.Trigger != "", h.HxTriggerString(props.Trigger)),
		h.If(props.Get != "", h.Get(props.Get)),
		h.If(props.Post != "", h.Post(props.Post)),
	)

	// If we only have an input with no additional elements, return it directly
	if props.Label == "" && props.Description == "" &&
		props.HelpText == nil && props.Error == "" &&
		props.LeadingIcon == nil && props.TrailingIcon == nil {
		return input
	}

	// Create wrapper with label, description, and input
	wrapperClasses := h.MergeClasses(
		"input-wrapper space-y-1.5",
		props.WrapClass,
	)

	children := make([]h.Ren, 0)

	// Add label if provided
	if props.Label != "" {
		labelClasses := h.MergeClasses(
			"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
			props.LabelClass,
		)
		children = append(children, h.Label(
			h.For(props.Id),
			h.Class(labelClasses),
			h.Text(props.Label),
		))
	}

	// Add description if provided
	if props.Description != "" {
		descriptionClasses := h.MergeClasses(
			"text-sm text-muted-foreground",
			props.DescriptionClass,
		)
		children = append(children, h.P(
			h.Class(descriptionClasses),
			h.Text(props.Description),
		))
	}

	// Create input container for icons
	// If we have icons, wrap the input in a relative container that matches input width
	if props.LeadingIcon != nil || props.TrailingIcon != nil {
		inputWrapperClasses := h.MergeClasses(
			"relative",
			h.Ternary(props.FullWidth, "w-full", "w-[320px]"), // Match input width
		)

		iconChildren := make([]h.Ren, 0)

		if props.LeadingIcon != nil {
			iconChildren = append(iconChildren, h.Div(
				h.Class("absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none"),
				props.LeadingIcon,
			))
		}

		iconChildren = append(iconChildren, input)

		if props.TrailingIcon != nil {
			iconChildren = append(iconChildren, h.Div(
				h.Class("absolute right-3 top-1/2 -translate-y-1/2"),
				props.TrailingIcon,
			))
		}

		children = append(children, h.Div(
			h.Class(inputWrapperClasses),
			h.Children(iconChildren...),
		))
	} else {
		children = append(children, input)
	}

	// Add help text if provided
	if props.HelpText != nil {
		children = append(children, h.Div(
			h.Class("text-sm text-muted-foreground mt-1"),
			props.HelpText,
		))
	}

	// Add error message if provided
	if props.Error != "" {
		errorClasses := h.MergeClasses(
			"text-sm font-medium text-destructive mt-1",
			props.ErrorClass,
		)
		children = append(children, h.P(
			h.Class(errorClasses),
			h.Text(props.Error),
		))
	}

	return h.Div(
		h.Class(wrapperClasses),
		h.Children(children...),
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

// TextareaProps extends InputProps for textarea-specific properties
type TextareaProps struct {
	InputProps
	Rows int
}

// Textarea creates a new textarea component with the provided props
func Textarea(props TextareaProps) *h.Element {
	baseClasses := "flex min-h-[80px] w-full rounded-md border border-input " +
		"bg-background px-3 py-2 text-sm ring-offset-background " +
		"placeholder:text-muted-foreground focus-visible:outline-none " +
		"focus-visible:ring-2 focus-visible:ring-gray-950 " +
		"dark:focus-visible:ring-gray-300 focus-visible:ring-offset-2 " +
		"disabled:cursor-not-allowed disabled:bg-background/50 outline-none"

	classes := h.MergeClasses(
		baseClasses,
		h.Ternary(props.Error != "", "border-destructive", ""),
		props.Class,
	)

	textarea := h.TextArea(
		h.Name(props.Name),
		h.Id(props.Id),
		h.Value(props.Value),
		h.Placeholder(props.Placeholder),
		h.Class(classes),
		h.If(props.Required, h.Required()),
		h.If(props.Disabled, h.Disabled()),
		h.If(props.ReadOnly, h.ReadOnly()),
		h.If(props.Rows > 0, h.Rows(props.Rows)),
	)

	// If we only have a textarea with no additional elements, return it directly
	if props.Label == "" && props.Description == "" &&
		props.HelpText != nil && props.Error == "" {
		return textarea
	}

	// Create wrapper with label, description, and textarea
	children := make([]h.Ren, 0)

	// Add label if provided
	if props.Label != "" {
		labelClasses := h.MergeClasses(
			"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
			props.LabelClass,
		)
		children = append(children, h.Label(
			h.For(props.Id),
			h.Class(labelClasses),
			h.Text(props.Label),
		))
	}

	// Add description if provided
	if props.Description != "" {
		descriptionClasses := h.MergeClasses(
			"text-sm text-muted-foreground",
			props.DescriptionClass,
		)
		children = append(children, h.P(
			h.Class(descriptionClasses),
			h.Text(props.Description),
		))
	}

	children = append(children, textarea)

	// Add help text if provided
	if props.HelpText != nil {
		children = append(children, h.P(
			h.Class("text-sm text-muted-foreground mt-1"),
			props.HelpText,
		))
	}

	// Add error message if provided
	if props.Error != "" {
		errorClasses := h.MergeClasses(
			"text-sm font-medium text-destructive mt-1",
			props.ErrorClass,
		)
		children = append(children, h.P(
			h.Class(errorClasses),
			h.Text(props.Error),
		))
	}

	return h.Div(
		h.Class("space-y-1.5"),
		h.Children(children...),
	)
}
