package resources

type RunType int

const (
	RunTypeUnknown RunType = iota
	RunTypeDockerBuild
	RunTypeDockerRegistry
)

type EmptyBuildMeta struct{}

type DockerBuildMeta struct {
	RepositoryUrl     string   `json:"repository_url"`
	Dockerfile        string   `json:"dockerfile"`
	GithubAccessToken string   `json:"github_access_token"`
	Tags              []string `json:"tags"`
}

type DockerRegistryMeta struct {
	Image string `json:"image"`
}

type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
