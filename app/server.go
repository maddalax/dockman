package app

import (
	"errors"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"paas/app/util/json2"
	"time"
)

type ServerPutOpts struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	LocalIpAddress  string    `json:"ip_address"`
	RemoteIpAddress string    `json:"remote_ip_address"`
	HostName        string    `json:"host_name"`
	LastSeen        time.Time `json:"last_seen"`
	Os              string    `json:"os"`
}

func ServerPut(locator *service.Locator, opts ServerPutOpts) error {
	client := service.Get[KvClient](locator)

	if opts.Id == "" {
		return errors.New("server id is required")
	}

	server, err := ServerGet(locator, opts.Id)

	// server exists already
	if err == nil {
		server.LocalIpAddress = opts.LocalIpAddress
		server.RemoteIpAddress = opts.RemoteIpAddress
		server.HostName = opts.HostName
		server.LastSeen = opts.LastSeen
		server.Os = opts.Os
		if opts.Name != "" {
			server.Name = opts.Name
		}
	} else {
		server = &Server{
			Id:              opts.Id,
			LocalIpAddress:  opts.LocalIpAddress,
			RemoteIpAddress: opts.RemoteIpAddress,
			Name:            opts.Name,
			HostName:        opts.HostName,
			LastSeen:        opts.LastSeen,
			Os:              opts.Os,
		}
	}

	bucket, err := client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "servers",
	})

	if err != nil {
		return err
	}

	err = client.PutJson(bucket, server.Id, server)

	if err != nil {
		return err
	}

	return nil
}

func ServerGet(locator *service.Locator, id string) (*Server, error) {
	client := service.Get[KvClient](locator)

	bucket, err := client.GetBucket("servers")

	if err != nil {
		return nil, err
	}

	server, err := bucket.Get(id)

	if err != nil || server == nil {
		return nil, err
	}

	return json2.Deserialize[Server](server.Value())
}

func ServerList(locator *service.Locator) ([]*Server, error) {
	client := service.Get[KvClient](locator)

	bucket, err := client.GetBucket("servers")

	if err != nil {
		return nil, err
	}

	var servers []*Server

	keys, err := bucket.ListKeys()

	if err != nil {
		return nil, err
	}

	for s := range keys.Keys() {
		server, err := ServerGet(locator, s)

		if err != nil {
			return nil, err
		}

		servers = append(servers, server)
	}

	return servers, nil
}
