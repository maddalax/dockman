package app

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"paas/app/logger"
	"sync"
	"time"
)

func (a *Agent) SubscribeToCommands() {
	_, err := a.kv.SubscribeSubject(context.Background(), "commands", func(msg *nats.Msg) {
		var wrapper struct {
			Command Command
			Id      string
		}

		buffer := bytes.NewBuffer(msg.Data)

		decoder := gob.NewDecoder(buffer)
		if err := decoder.Decode(&wrapper); err != nil {
			logger.Error("Failed to decode command: %s", err)
			return
		}

		wrapper.Command.Execute(a)

		response := wrapper.Command.GetResponse()

		serialized, err := GobSerializeResponse(response)

		if err != nil {
			logger.Error("Failed to serialize response", err)
			return
		}

		bucket := a.commandResponseBucket

		_, err = bucket.Put(wrapper.Id, serialized.Bytes())

		if err != nil {
			logger.Error("Failed to put response", err)
		}
	})

	if err != nil {
		return
	}
}

type SendCommandResponse[T any] struct {
	Response      T
	ServerDetails ServerDetails
}

type SendCommandOpts struct {
	ExpectedResponses int
	Command           Command
	Timeout           time.Duration
}

func SendCommand[T any](locator *service.Locator, opts SendCommandOpts) ([]*SendCommandResponse[T], error) {
	agent := AgentFromLocator(locator)

	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	var responses = make([]*SendCommandResponse[T], 0)

	buffer := bytes.Buffer{}

	encoder := gob.NewEncoder(&buffer)
	cmd := NewCommand(opts.Command)

	err := encoder.Encode(cmd)

	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(opts.Timeout)
		defer ticker.Stop()

		watcher, err := agent.commandResponseBucket.Watch(cmd.Id)

		if err != nil {
			return
		}

		defer watcher.Stop()

		for {
			select {
			case <-ticker.C:
				return
			case c := <-watcher.Updates():

				if c == nil {
					continue
				}

				decoder := gob.NewDecoder(bytes.NewBuffer(c.Value()))

				var responseWrapper ResponseWrapper[any]
				err := decoder.Decode(&responseWrapper)

				if err != nil {
					logger.Error("Failed to decode response", err)
					return
				}
				details := responseWrapper.ServerDetails
				cast, ok := responseWrapper.Response.(*T)
				if !ok {
					logger.Error("unable to cast command response", errors.New("failed to cast response"))
					return
				}
				responses = append(responses, &SendCommandResponse[T]{
					Response:      *cast,
					ServerDetails: details,
				})

				// If we have received all the responses we were expecting, return
				if len(responses) == opts.ExpectedResponses {
					return
				}
			}
		}

	}()

	_, err = agent.commandWriter.Write(buffer.Bytes())

	wg.Wait()

	return responses, err
}
