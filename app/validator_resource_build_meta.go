package app

import (
	"bufio"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"os"
	"path/filepath"
	"strings"
)

type BuildMetaValidator struct {
	Meta BuildMeta
}

func (v BuildMetaValidator) Validate() error {
	var validators []Validator

	switch m := v.Meta.(type) {
	case *DockerBuildMeta:
		validators = []Validator{
			GithubRepositoryValidator{
				RepositoryUrl: m.RepositoryUrl,
				AccessToken:   m.GithubAccessToken,
				Dockerfile:    m.Dockerfile,
			},
		}
	}

	for _, validator := range validators {
		err := validator.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type GithubRepositoryValidator struct {
	RepositoryUrl string
	AccessToken   string
	Dockerfile    string
}

func (v GithubRepositoryValidator) Validate() error {
	if v.RepositoryUrl == "" {
		return errors.New("git repository url is required")
	}

	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{v.RepositoryUrl},
	})

	opts := &git.ListOptions{}

	if v.AccessToken != "" {
		opts.Auth = &http.BasicAuth{
			Username: "paas",
			Password: v.AccessToken,
		}
	}

	_, err := rem.List(opts)

	if err != nil {
		if err.Error() == "authentication required" {
			return errors.New("repository is not accessible, please ensure you have provided a personal access token with 'Contents' permission")
		}
		if err.Error() == "repository not found" {
			return errors.New("repository not found, please ensure the url is correct")
		}
		return err
	}

	if v.Dockerfile != "" {
		clone, err := Clone(CloneRequest{
			Meta: &DockerBuildMeta{
				RepositoryUrl:     v.RepositoryUrl,
				Dockerfile:        v.Dockerfile,
				GithubAccessToken: v.AccessToken,
			},
			Progress: os.Stdout,
		})

		if err != nil {
			return err
		}

		validator := ValidDockerFileValidator{
			Dockerfile:    v.Dockerfile,
			RepositoryDir: clone.Directory,
		}
		return validator.Validate()
	}

	return nil
}

type ValidDockerFileValidator struct {
	Dockerfile    string
	RepositoryDir string
}

func (v ValidDockerFileValidator) Validate() error {
	dockerfilePath := filepath.Join(v.RepositoryDir, v.Dockerfile)
	_, err := os.Lstat(dockerfilePath)

	if err != nil {
		return errors.New("dockerfile not found, please ensure the path is correct and is relative from the repository root")
	}

	// validate it's a valid dockerfile with a quick check
	file, err := os.Open(dockerfilePath)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// all good
		if strings.HasPrefix(strings.ToLower(line), "from") {
			return nil
		}
	}

	return errors.New("found the specified Dockerfile but it didn't appear to be a valid Dockerfile")
}
