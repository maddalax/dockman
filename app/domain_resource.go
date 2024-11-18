package app

import (
	"encoding/json"
	"fmt"
	"paas/json2"
	"slices"
)

type Resource struct {
	Id                 string            `json:"id"`
	Name               string            `json:"name"`
	Environment        string            `json:"environment"`
	RunType            RunType           `json:"run_type"`
	InstancesPerServer int               `json:"instances_per_server"`
	BuildMeta          BuildMeta         `json:"build_meta"`
	Env                map[string]string `json:"env"`
	RunStatus          RunStatus         `json:"run_status"`
}

func (resource *Resource) MarshalJSON() ([]byte, error) {

	if resource.Env == nil {
		resource.Env = make(map[string]string)
	}

	if resource.InstancesPerServer == 0 {
		resource.InstancesPerServer = 1
	}

	buildMeta := json2.SerializeOrEmpty(resource.BuildMeta)
	return json.Marshal(map[string]interface{}{
		"id":                   resource.Id,
		"name":                 resource.Name,
		"environment":          resource.Environment,
		"run_type":             resource.RunType,
		"instances_per_server": resource.InstancesPerServer,
		"build_meta":           json.RawMessage(buildMeta),
		"env":                  resource.Env,
		"run_status":           resource.RunStatus,
	})
}

func (resource *Resource) UnmarshalJSON(data []byte) error {
	temp := make(map[string]interface{})
	err := json.Unmarshal(data, &temp)

	if err != nil {
		return err
	}

	resource.Id = temp["id"].(string)
	resource.Name = temp["name"].(string)
	resource.Environment = temp["environment"].(string)
	resource.InstancesPerServer = int(temp["instances_per_server"].(float64))
	resource.RunStatus = RunStatus(temp["run_status"].(float64))
	resource.RunType = RunType(temp["run_type"].(float64))

	env, ok := temp["env"].(map[string]interface{})

	if ok {
		for k, v := range env {
			resource.Env[k] = v.(string)
		}
	}

	buildMeta := temp["build_meta"].(map[string]interface{})
	switch resource.RunType {
	case RunTypeDockerBuild:
		resource.BuildMeta = &DockerBuildMeta{}
		serialized := json2.SerializeOrEmpty(buildMeta)
		err = json.Unmarshal(serialized, &resource.BuildMeta)
	case RunTypeDockerRegistry:
		resource.BuildMeta = &DockerRegistryMeta{
			Image: buildMeta["image"].(string),
		}
	default:
		resource.BuildMeta = &EmptyBuildMeta{}
	}

	return nil
}

type RunStatus = int

const (
	RunStatusUnknown RunStatus = iota
	RunStatusNotRunning
	RunStatusRunning
	RunStatusPartiallyRunning
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

type BuildMeta interface {
	ValidatePatch(other BuildMeta) error
}

type ValidateBuildMetaPatchResponse struct {
	DidChange bool
}

type EmptyBuildMeta struct{}

func (bm *EmptyBuildMeta) ValidatePatch(other BuildMeta) error {
	return nil
}

type DockerBuildMeta struct {
	RepositoryUrl     string   `json:"repository_url"`
	Dockerfile        string   `json:"dockerfile"`
	GithubAccessToken string   `json:"github_access_token"`
	Tags              []string `json:"tags"`
	ExposedPort       int      `json:"exposed_port"`
}

func (bm *DockerBuildMeta) ValidatePatch(other BuildMeta) error {
	b2, ok := other.(*DockerBuildMeta)

	if !ok {
		return fmt.Errorf("invalid build meta type")
	}
	response := &ValidateBuildMetaPatchResponse{}

	response.DidChange = bm.RepositoryUrl == b2.RepositoryUrl &&
		bm.Dockerfile == b2.Dockerfile &&
		bm.GithubAccessToken == b2.GithubAccessToken &&
		slices.Equal(bm.Tags, b2.Tags) &&
		bm.ExposedPort == b2.ExposedPort

	// repository access may have changed, re-validate it
	if bm.RepositoryUrl != b2.RepositoryUrl || bm.GithubAccessToken != b2.GithubAccessToken || bm.Dockerfile != b2.Dockerfile {
		validator := BuildMetaValidator{
			Meta: b2,
		}
		err := validator.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type DockerRegistryMeta struct {
	Image string `json:"image"`
}

func (bm *DockerRegistryMeta) ValidatePatch(other BuildMeta) error {
	b2, ok := other.(*DockerRegistryMeta)

	if !ok {
		return fmt.Errorf("build meta type mismatch")
	}

	response := &ValidateBuildMetaPatchResponse{}
	response.DidChange = bm.Image == b2.Image

	return nil
}

type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
