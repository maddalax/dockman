package main

import (
	"dockside/app"
	"dockside/app/logger"
	"log"
	"strings"

	"github.com/kardianos/service"
	service2 "github.com/maddalax/htmgo/framework/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
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

func (p *program) Stop(s service.Service) error {

	// Clean up here.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "dockside-agent",
		DisplayName: "Dockside Agent",
		Description: "Dockside Agent",
	}

	program := &program{}
	s, err := service.New(program, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Install()

	if err != nil {
		if strings.HasPrefix(err.Error(), "Init already exists") {
			// do nothing
		} else {
			log.Fatal(err)
		}
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
