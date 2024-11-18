package docker

import (
	"context"
	"errors"
	"github.com/docker/docker/client"
	"sync"
)

type Client struct {
	cli *client.Client
}

var _client *Client
var syncOnce sync.Once

func Connect() (*Client, error) {
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
			_client = &Client{
				cli: cli,
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
