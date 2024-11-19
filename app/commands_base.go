package app

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
	"os"
)

type Command interface {
	Execute(agent *Agent)
	GetResponse() any
}

func NewCommand(command Command) CommandWrapper[Command] {
	wrapper := CommandWrapper[Command]{Command: command}
	wrapper.Id = uuid.NewString()
	return wrapper
}

type CommandWrapper[T Command] struct {
	Command T
	Id      string
}

type ResponseWrapper[T any] struct {
	Response      T
	ServerDetails ServerDetails
}

type ServerDetails struct {
	Hostname string
}

func GobSerializeResponse[T any](data T) (bytes.Buffer, error) {
	hostName, _ := os.Hostname()

	wrapper := ResponseWrapper[any]{
		Response: data,
		ServerDetails: ServerDetails{
			Hostname: hostName,
		},
	}

	gob.Register(data)
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wrapper)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buffer, nil
}
