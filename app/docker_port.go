package app

import (
	"fmt"
	"net"
)

func FindOpenPort(startPort int) (int, error) {
	for port := startPort; port <= 65535; port++ {
		address := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no open ports found starting from %d", startPort)
}
