package app

import (
	"dockside/app/logger"
	"github.com/nats-io/nats-server/v2/server"
	"time"
)

func MustStartNats() *server.Server {
	logger.Info("Starting NATS server")

	opts := &server.Options{
		Port:      4222,
		JetStream: true,
		StoreDir:  "./data",
	}

	natsServer, err := server.NewServer(opts)
	if err != nil {
		panic(err)
	}

	go natsServer.Start()

	if !natsServer.ReadyForConnections(5 * time.Second) {
		panic("Failed to start NATS server after 5 seconds")
	}

	logger.Info("NATS server started")

	return natsServer
}
