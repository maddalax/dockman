package main

import (
	"dockside/app"
	"dockside/app/logger"
	"github.com/maddalax/htmgo/framework/service"
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
		logger.Error("Failed to start fluentd container, unable to stream logs", err)
	}

	go func() {
		err := fluentd.StreamLogs()
		if err != nil {
			logger.Error("Failed to stream logs from fluentd", err)
		}
	}()

	agent.Run()
}
