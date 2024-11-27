package app

import (
	"dockside/app/logger"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func FindOpenPort() (int, error) {
	minp := 1024
	maxp := 49151
	attempts := 0
	for {
		attempts++
		if attempts > 100 {
			break
		}
		randomPort := rand.Intn(maxp-minp+1) + minp
		logger.InfoWithFields("Trying to find open port", map[string]any{
			"port": randomPort,
		})
		address := fmt.Sprintf(":%d", randomPort)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			listener.Close()
			return randomPort, nil
		}
		time.Sleep(time.Millisecond * 10)
	}

	return 0, fmt.Errorf("no open ports found starting from %d to %d", minp, maxp)
}
