package app

import (
	"encoding/gob"
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"time"
)

type Agent struct {
	setup                 bool
	locator               *service.Locator
	kv                    *KvClient
	commandWriter         *EphemeralNatsWriter
	commandResponseBucket nats.KeyValue
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
			Port: 4222,
		})
		if err != nil {
			panic(err)
		}
		return client
	})

	a.kv = KvFromLocator(a.locator)
	a.commandWriter = a.kv.NewEphemeralNatsWriter("commands")

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

	for {
		fmt.Printf("Agent is running\n")
		time.Sleep(time.Second * 5)
	}
}
