package main

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/app"
)

func main() {
	locator := service.NewLocator()
	agent := app.NewAgent(locator)
	err := agent.Setup()
	if err != nil {
		panic(err)
	}

	fluentd := NewFluentdManager(agent)
	err = fluentd.StartContainer()

	if err != nil {
		panic(err)
	}

	go fluentd.StreamLogs()

	agent.Run()
}
