package internal

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
)

type CloneRepoResult struct {
	Directory string
}

type CloneRepoRequest struct {
	Progress sideband.Progress
}

func (bm *DockerBuildMeta) CloneRepo(request CloneRepoRequest) (*CloneRepoResult, error) {
	tempDir, err := os.MkdirTemp("", "repo-clone-*")

	if err != nil {
		return nil, err
	}

	os.Chmod(tempDir, 0700)

	opts := &git.CloneOptions{}

	if bm.GithubAccessToken != "" {
		opts.Auth = &http.BasicAuth{
			Username: "paas",
			Password: bm.GithubAccessToken,
		}
	}

	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      bm.RepositoryUrl,
		Auth:     opts.Auth,
		Progress: request.Progress,
		Depth:    1,
	})

	if err != nil {
		return nil, err
	}

	return &CloneRepoResult{
		Directory: tempDir,
	}, nil
}
