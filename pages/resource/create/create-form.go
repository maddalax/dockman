package create

import (
	"dockside/app/ui"
	"dockside/app/ui/icons"
	"github.com/maddalax/htmgo/framework/h"
)

func DeploymentChoiceSelector() *h.Element {
	return h.Div(
		ui.FieldLabel("Deployment Type"),
		h.FieldSet(
			h.Class("flex flex-wrap gap-2 mt-1"),
			h.Tag(
				"legend",
				h.Class("sr-only"),
			),
			h.Div(
				h.Class("max-w-[250px]"),
				DockerFileChoice(),
			),
			h.Div(
				h.Class("max-w-[250px]"),
				DockerRegistryChoice(),
			),
		),
	)
}

func DockerFileChoice() *h.Element {
	t := "dockerfile"
	return ui.ChoiceCard(ui.ChoiceCardProps{
		Title:          "Dockerfile",
		Description:    "Build your application from a specified Dockerfile",
		Icon:           icons.DockerFileIconBlack(),
		InputName:      "deployment-type",
		InputValue:     t,
		Id:             t,
		DefaultChecked: false,
		InputProps: h.GetPartialWithQs(
			AdditionalCreateResourceFields,
			h.NewQs("deployment_type", t),
			"change",
		),
	})
}

func DockerRegistryChoice() *h.Element {
	t := "docker_registry"
	return ui.ChoiceCard(ui.ChoiceCardProps{
		Title:          "Docker Registry",
		Description:    "Run your application from an existing Docker image",
		Icon:           icons.DockerIconBlack(),
		InputName:      "deployment-type",
		InputValue:     t,
		Id:             t,
		DefaultChecked: false,
		InputProps: h.GetPartialWithQs(
			AdditionalCreateResourceFields,
			h.NewQs("deployment_type", t),
			"change",
		),
	})
}
