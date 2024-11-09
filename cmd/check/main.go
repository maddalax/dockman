package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	ip := "192.168.4.70" // Replace with the IP address of the device you want to check
	port := 6380

	serviceInfo := checkService(ip, port)
	if serviceInfo != "" {
		fmt.Printf("Service running on %s:%d - %s\n", ip, port, serviceInfo)
	} else {
		fmt.Printf("Unable to identify service on %s:%d\n", ip, port)
	}
}

// checkService tries to identify the service running on a specific port
func checkService(ip string, port int) string {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return ""
	}
	defer conn.Close()

	// Send a simple Redis PING command
	fmt.Fprintf(conn, "*1\r\n$4\r\nPING\r\n")

	// Read the response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	// Check if the response matches a known Redis response
	if strings.HasPrefix(response, "+PONG") {
		return "Redis"
	}

	return "Unknown service"
}
