package resources

type RunType int

const (
	RunTypeDockerBuild RunType = iota
	RunTypeDockerRegistry
)

type DockerBuildMeta struct {
	Dockerfile string   `json:"dockerfile"`
	Tags       []string `json:"tags"`
}

type DockerRegistryMeta struct {
	Image string `json:"image"`
}

type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Resource struct {
	Name        string  `json:"name"`
	Environment string  `json:"environment"`
	RunType     RunType `json:"run_type"`
	BuildMeta   any     `json:"build_meta"`
	Env         []Env   `json:"env"`
}
