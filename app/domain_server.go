package app

import "time"

type Server struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	LocalIpAddress  string    `json:"ip_address"`
	RemoteIpAddress string    `json:"remote_ip_address"`
	HostName        string    `json:"host_name"`
	LastSeen        time.Time `json:"last_seen"`
	Os              string    `json:"os"`
}
