package app

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
)

type Command interface {
	Execute(agent *Agent)
	GetResponse() any
	Name() string
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
	ServerDetails Server
}

func GobSerializeResponse[T any](data T) (bytes.Buffer, error) {
	wrapper := ResponseWrapper[any]{
		Response: data,
	}

	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wrapper)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buffer, nil
}
