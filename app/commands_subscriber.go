package app

import (
	"bytes"
	"context"
	"dockside/app/logger"
	"encoding/gob"
	"errors"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"sync"
	"time"
)

func (a *Agent) SubscribeToCommands() {
	logger.InfoWithFields("Subscribing to commands", map[string]any{
		"subject":  a.commandStreamName,
		"serverId": a.serverId,
	})

	_, err := a.kv.SubscribeSubject(context.Background(), a.commandStreamName, func(msg *nats.Msg) {
		logger.InfoWithFields("Received command", map[string]any{
			"subject": msg.Subject,
			"size":    msg.Size(),
		})

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

		logger.InfoWithFields("executing command", map[string]any{
			"command": wrapper.Command,
			"id":      wrapper.Id,
		})
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
	ServerIds         []string
	Timeout           time.Duration
}

func SendCommandForResource[T any](locator *service.Locator, resourceId string, opts SendCommandOpts) ([]*SendCommandResponse[T], error) {
	agent := AgentFromLocator(locator)
	serverIds, err := ResourceGetServerIds(agent.locator, resourceId)
	if err != nil {
		return nil, err
	}
	opts.ServerIds = serverIds
	responses := make([]*SendCommandResponse[T], 0)

	for _, id := range opts.ServerIds {
		result, err := SendCommand[T](locator, id, opts)
		if err != nil {
			logger.ErrorWithFields("Failed to send command", err, map[string]any{
				"server_id": id,
			})
			continue
		}
		if result != nil {
			responses = append(responses, result)
		}
	}

	return responses, nil
}

func SendCommand[T any](locator *service.Locator, serverId string, opts SendCommandOpts) (*SendCommandResponse[T], error) {
	agent := AgentFromLocator(locator)

	if serverId == "" {
		return nil, errors.New("server id must be provided")
	}

	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	var response = new(SendCommandResponse[T])

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
				response = &SendCommandResponse[T]{
					Response:      *cast,
					ServerDetails: details,
				}
			}
		}

	}()

	logger.InfoWithFields("sending command", map[string]any{
		"command":    opts.Command,
		"server_ids": opts.ServerIds,
	})

	subjectName := agent.CommandStreamName(serverId)

	writer := agent.kv.NewEphemeralNatsWriter(subjectName)
	defer writer.Close()

	_, err = writer.Write(buffer.Bytes())

	wg.Wait()

	return response, err
}
