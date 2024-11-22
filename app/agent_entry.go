package app

import (
	"dockside/app/logger"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"time"
)

type Agent struct {
	setup             bool
	locator           *service.Locator
	registry          *ServiceRegistry
	commandStreamName string
	serverId          string
}

func (a *Agent) GetLocator() *service.Locator {
	return a.locator
}

func (a *Agent) GetServerId() string {
	return a.serverId
}

func NewAgent(locator *service.Locator) *Agent {
	return &Agent{
		locator:  locator,
		registry: GetServiceRegistry(locator),
	}
}

func AgentFromLocator(locator *service.Locator) *Agent {
	return service.Get[Agent](locator)
}

func (a *Agent) GetCommandResponseBucket() (nats.KeyValue, error) {
	bucket, err := a.registry.KvClient().GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "command_responses",
		TTL:    time.Hour,
	})
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (a *Agent) Setup() error {

	a.RegisterGobTypes()

	a.registry = GetServiceRegistry(a.locator)

	serverId := a.registry.GetServerConfigManager().GetConfig("server_id")

	// no server id set, generate one
	if serverId == "" {
		a.serverId = uuid.NewString()
		a.registry.GetServerConfigManager().WriteConfig("server_id", a.serverId)
	} else {
		a.serverId = serverId
	}

	a.commandStreamName = a.CommandStreamName(a.serverId)

	return nil
}

func (a *Agent) CommandStreamName(serverId string) string {
	return "commands-" + serverId
}

func (a *Agent) RegisterGobTypes() {
	gob.Register(&RunResourceCommand{})
	gob.Register(&RunResourceResponse{})
	gob.Register(&StopResourceCommand{})
	gob.Register(&StopResourceResponse{})
	gob.Register(&PingCommand{})
	gob.Register(&PingResponse{})
	gob.Register(&SetServerConfigCommand{})
	gob.Register(&GetContainerCommand{})
	gob.Register(&GetContainerResponse{})
}

func (a *Agent) Run() {
	if !a.setup {
		a.Setup()
		a.setup = true
	}

	err := a.registry.KvClient().Ping()

	if err != nil {
		panic(err)
	}

	a.SubscribeToCommands()
	a.RegisterMonitor()

	go a.registry.GetJobRunner().Start()

	for {
		logger.Info("Agent is running")
		time.Sleep(time.Second * 5)
	}
}
