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

func (server *Server) IsAccessible() bool {
	now := time.Now()
	// has sent an update in the last 10 seconds
	return now.Sub(server.LastSeen) < time.Second*10
}

func (server *Server) FormattedName() string {
	if server.Name != "" {
		return server.Name
	}
	return server.HostName
}
