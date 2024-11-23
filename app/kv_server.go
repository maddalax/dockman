package app

import (
	"dockside/app/logger"
	"github.com/nats-io/nats-server/v2/server"
	"runtime"
	"time"
)

func MustStartNats() *server.Server {
	logger.Info("Starting NATS server")

	storeDir := "/data/dockside"

	if runtime.GOOS == "windows" {
		storeDir = "C:/data/dockside"
	}

	if runtime.GOOS == "darwin" {
		storeDir = "~/.dockside/data"
	}

	opts := &server.Options{
		Port:      4222,
		JetStream: true,
		StoreDir:  storeDir,
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
