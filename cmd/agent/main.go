package main

import (
	"github.com/maddalax/htmgo/framework/service"
	"paas/internal"
)

func main() {
	locator := service.NewLocator()
	agent := internal.NewAgent(locator)
	agent.Setup()
	agent.Run()
}
