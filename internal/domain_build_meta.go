package internal

import (
	"fmt"
	"slices"
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
