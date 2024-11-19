package app

import (
	"fmt"
	"github.com/nats-io/nats-server/v2/server"
	"time"
)

func StartNatsServer() (*server.Server, error) {
	// ResourceCreate a NATS server configuration
	opts := &server.Options{
		Port:      4222, // You can choose a different port if needed
		JetStream: true,
		StoreDir:  "./data",
	}

	// ResourceCreate a new NATS server instance
	natsServer, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}

	// Start the NATS server in a goroutine
	go natsServer.Start()

	// Check if the server is ready
	if !natsServer.ReadyForConnections(5 * time.Second) {
		return nil, fmt.Errorf("NATS server did not start in time")
	}

	return natsServer, nil
}
