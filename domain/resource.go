package domain

import "fmt"

type Resource struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	RunType     RunType           `json:"run_type"`
	BuildMeta   any               `json:"build_meta"`
	Env         map[string]string `json:"env"`
	RunStatus   RunStatus         `json:"run_status"`
}

func (resource *Resource) BucketKey() string {
	return fmt.Sprintf("resources-%s", resource.Id)
}

type RunStatus = int

const (
	RunStatusUnknown RunStatus = iota
	RunStatusNotRunning
	RunStatusRunning
	RunStatusErrored
)

func NewResource(id string) *Resource {
	resource := &Resource{
		Id: id,
	}
	return resource
}

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
