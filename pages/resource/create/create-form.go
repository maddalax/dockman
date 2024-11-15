package create

import (
	"github.com/maddalax/htmgo/framework/h"
	"paas/ui"
	"paas/ui/icons"
)

func DeploymentChoiceSelector() *h.Element {
	return h.Div(
		h.H2F("Deployment Type"),
		h.FieldSet(
			h.Class("flex flex-wrap mt-1"),
			h.Tag("legend", h.Class("sr-only")),
			DockerFileChoice(),
			DockerRegistryChoice(),
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
		InputProps:     h.GetPartialWithQs(AdditionalCreateResourceFields, h.NewQs("deployment_type", t), "change"),
	})
}

func DockerRegistryChoice() *h.Element {
	t := "docker_registry"
	return ui.ChoiceCard(ui.ChoiceCardProps{
		Title:          "Docker Registry",
		Description:    "Run your application from an existing Docker image in a registry",
		Icon:           icons.DockerIconBlack(),
		InputName:      "deployment-type",
		InputValue:     t,
		Id:             t,
		DefaultChecked: false,
		InputProps:     h.GetPartialWithQs(AdditionalCreateResourceFields, h.NewQs("deployment_type", t), "change"),
	})
}
