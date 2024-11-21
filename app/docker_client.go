package app

import (
	"context"
	"errors"
	"github.com/docker/docker/client"
	"github.com/maddalax/htmgo/framework/service"
	"sync"
)

type DockerClient struct {
	cli     *client.Client
	locator *service.Locator
}

var _client *DockerClient
var syncOnce sync.Once

func DockerConnect(locator *service.Locator) (*DockerClient, error) {
	var e error
	syncOnce.Do(func() {
		env := client.FromEnv
		cli, err := client.NewClientWithOpts(env,
			client.WithAPIVersionNegotiation(),
		)
		if err != nil {
			e = err
		}
		if e == nil {
			_client = &DockerClient{
				cli:     cli,
				locator: locator,
			}
		}
	})

	if _client == nil {
		return nil, errors.New("failed to connect to docker")
	}

	_, err := _client.cli.Ping(context.Background())

	if err != nil {
		return nil, err
	}

	return _client, e
}
