package main

import (
	"dockside/app"
	"dockside/app/logger"
	service2 "github.com/maddalax/htmgo/framework/service"
)

func main() {
	locator := service2.NewLocator()
	registry := app.CreateServiceRegistry(locator)

	registry.RegisterAgentStartupServices()

	agent := registry.GetAgent()

	fluentd := NewFluentdManager(agent)
	err := fluentd.StartContainer()

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
