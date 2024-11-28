package app

import (
	"bytes"
	"dockman/app/logger"
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

	_, err := a.registry.KvClient().SubscribeSubjectForever(a.commandStreamName, func(msg *nats.Msg) {
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
			logger.Error("Failed to decode command", err)
			return
		}

		logger.InfoWithFields("executing command", map[string]any{
			"command": wrapper.Command.Name(),
			"id":      wrapper.Id,
		})
		wrapper.Command.Execute(a)

		response := wrapper.Command.GetResponse()

		serialized, err := GobSerializeResponse(response)

		if err != nil {
			logger.Error("Failed to serialize response", err)
			return
		}

		bucket, err := a.GetCommandResponseBucket()

		if err != nil {
			logger.Error("Failed to get command response bucket", err)
			return
		}

		_, err = bucket.Put(wrapper.Id, serialized.Bytes())

		if err != nil {
			logger.Error("Failed to put response", err)
		}
	})

	if err != nil {
		logger.ErrorWithFields("Failed to subscribe to commands", err, map[string]any{
			"subject": a.commandStreamName,
		})
	}
}

type SendCommandResponse[T any] struct {
	Response      T
	ServerDetails *Server
	SendError     error
}

type SendCommandOpts struct {
	ExpectedResponses int
	Command           Command
	Timeout           time.Duration
	PingFirst         bool
}

func SendCommandForResource[T any](locator *service.Locator, resourceId string, opts SendCommandOpts) ([]*SendCommandResponse[T], error) {
	agent := AgentFromLocator(locator)
	serverIds, err := ResourceGetServerIds(agent.locator, resourceId)
	if err != nil {
		return nil, err
	}
	responses := make([]*SendCommandResponse[T], 0)

	wg := sync.WaitGroup{}
	for _, id := range serverIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			opts.PingFirst = true
			result, err := SendCommand[T](locator, id, opts)
			if err != nil {
				logger.ErrorWithFields("Failed to send command", err, map[string]any{
					"server_id": id,
				})
				return
			}
			if result != nil {
				responses = append(responses, result)
			}
		}()
	}

	wg.Wait()

	return responses, nil
}

func SendPingToServer(locator *service.Locator, serverId string) bool {
	res, err := SendCommand[PingResponse](locator, serverId, SendCommandOpts{
		Command:   &PingCommand{},
		Timeout:   time.Second * 5,
		PingFirst: false,
	})
	if err != nil {
		return false
	}
	return res.Response.Message == "pong"
}

func SendCommand[T any](locator *service.Locator, serverId string, opts SendCommandOpts) (*SendCommandResponse[T], error) {
	agent := AgentFromLocator(locator)

	if serverId == "" {
		return nil, errors.New("server id must be provided")
	}

	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	var response = &SendCommandResponse[T]{
		ServerDetails: &Server{},
	}

	server, err := ServerGet(locator, serverId)

	if err == nil {
		response.ServerDetails = server
	}

	if opts.PingFirst {
		isAccessible := SendPingToServer(locator, serverId)

		if !isAccessible {
			response.SendError = errors.New("server is not accessible")
			return response, nil
		}
	}

	buffer := bytes.Buffer{}

	encoder := gob.NewEncoder(&buffer)
	cmd := NewCommand(opts.Command)

	err = encoder.Encode(cmd)

	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(opts.Timeout)
		defer ticker.Stop()

		bucket, err := agent.GetCommandResponseBucket()

		if err != nil {
			response.SendError = err
			return
		}

		watcher, err := bucket.Watch(cmd.Id)

		if err != nil {
			return
		}

		defer watcher.Stop()

		for {
			select {
			case <-ticker.C:
				// no response before the timeout
				response.SendError = errors.New("no response from server")
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

				cast, ok := responseWrapper.Response.(*T)
				if !ok {
					logger.Error("unable to cast command response", errors.New("failed to cast response"))
					return
				} else {
					newResponse := &SendCommandResponse[T]{
						Response: *cast,
					}
					response.Response = newResponse.Response
					return
				}
			}
		}

	}()

	subjectName := agent.CommandStreamName(serverId)

	logger.InfoWithFields("sending command", map[string]any{
		"command":   opts.Command.Name(),
		"server_id": serverId,
		"stream":    subjectName,
	})

	writer := agent.registry.KvClient().NewEphemeralNatsWriter(subjectName)
	defer writer.Close()

	_, err = writer.Write(buffer.Bytes())

	if err != nil {
		return nil, err
	}

	wg.Wait()

	return response, err
}
