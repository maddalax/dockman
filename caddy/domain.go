package caddy

type Config struct {
	Apps Apps `json:"apps"`
}

type Apps struct {
	HTTP HTTP `json:"http"`
}

type HTTP struct {
	Servers map[string]Server `json:"servers"`
}

type Server struct {
	Listen []string `json:"listen"`
	Routes []Route  `json:"routes"`
}

type Route struct {
	Handle   []Handler `json:"handle"`
	Match    []Match   `json:"match"`
	Terminal bool      `json:"terminal"`
}

type Handler struct {
	Handler   string     `json:"handler"`
	Routes    []SubRoute `json:"routes,omitempty"`
	Upstreams []Upstream `json:"upstreams,omitempty"`
}

type SubRoute struct {
	Handle []Handler `json:"handle"`
}

type Upstream struct {
	Dial string `json:"dial"`
}

type Match struct {
	Host []string `json:"host"`
}
