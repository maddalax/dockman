package app

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
)

type CloneResult struct {
	Directory string
}

type CloneRequest struct {
	Meta     *DockerBuildMeta
	Progress sideband.Progress
}

func Clone(request CloneRequest) (*CloneResult, error) {
	tempDir, err := os.MkdirTemp("", "repo-clone-*")

	if err != nil {
		return nil, err
	}

	os.Chmod(tempDir, 0700)

	opts := &git.CloneOptions{}

	if request.Meta.GithubAccessToken != "" {
		opts.Auth = &http.BasicAuth{
			Username: "paas",
			Password: request.Meta.GithubAccessToken,
		}
	}

	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      request.Meta.RepositoryUrl,
		Auth:     opts.Auth,
		Progress: request.Progress,
		Depth:    1,
	})

	if err != nil {
		return nil, err
	}

	return &CloneResult{
		Directory: tempDir,
	}, nil
}
