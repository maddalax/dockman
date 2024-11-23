package app

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
)

type CloneRepoResult struct {
	Directory string
	Commit    string
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
			Username: "dockside",
			Password: bm.GithubAccessToken,
		}
	}

	repo, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      bm.RepositoryUrl,
		Auth:     opts.Auth,
		Progress: request.Progress,
		Depth:    1,
	})

	if err != nil {
		return nil, err
	}

	commitHash := ""

	ref, err := repo.Head()
	if err == nil {
		commit, err := repo.CommitObject(ref.Hash())
		if err == nil {
			commitHash = commit.Hash.String()
		}
	}

	return &CloneRepoResult{
		Directory: tempDir,
		Commit:    commitHash,
	}, nil
}
