package app

import (
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"os"
	"paas/app/logger"
	"time"
)

type Agent struct {
	setup                 bool
	locator               *service.Locator
	kv                    *KvClient
	commandWriter         *EphemeralNatsWriter
	serverConfigManager   *ServerConfigManager
	commandResponseBucket nats.KeyValue
	serverId              string
}

func NewAgent(locator *service.Locator) *Agent {
	return &Agent{
		locator: locator,
	}
}

func AgentFromLocator(locator *service.Locator) *Agent {
	return service.Get[Agent](locator)
}

func (a *Agent) Setup() error {

	a.RegisterGobTypes()

	service.Set[KvClient](a.locator, service.Singleton, func() *KvClient {
		client, err := NatsConnect(NatsConnectOptions{
			Host: os.Getenv("NATS_HOST"),
			Port: 4222,
		})
		if err != nil {
			panic(err)
		}
		return client
	})

	service.Set[ServerConfigManager](a.locator, service.Singleton, func() *ServerConfigManager {
		return NewServerConfigManager()
	})

	a.kv = KvFromLocator(a.locator)
	a.commandWriter = a.kv.NewEphemeralNatsWriter("commands")
	a.serverConfigManager = service.Get[ServerConfigManager](a.locator)

	bucket, err := a.kv.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "command_responses",
		TTL:    time.Hour,
	})

	if err != nil {
		return err
	}

	a.commandResponseBucket = bucket

	service.Set(a.locator, service.Singleton, func() *Agent {
		return a
	})

	serverId := a.serverConfigManager.GetConfig("server_id")

	// no server id set, generate one
	if serverId == "" {
		a.serverId = uuid.NewString()
		a.serverConfigManager.WriteConfig("server_id", a.serverId)
	} else {
		a.serverId = serverId
	}

	return nil
}

func (a *Agent) RegisterGobTypes() {
	gob.Register(&RunResourceCommand{})
	gob.Register(&RunResourceResponse{})
	gob.Register(&StopResourceCommand{})
	gob.Register(&StopResourceResponse{})
}

func (a *Agent) Run() {
	if !a.setup {
		a.Setup()
		a.setup = true
	}

	err := a.kv.Ping()

	if err != nil {
		panic(err)
	}

	a.SubscribeToCommands()

	go a.StartServerMonitor()

	for {
		logger.Debug("Agent is running")
		time.Sleep(time.Second * 5)
	}
}
