package app

import (
	"dockside/app/logger"
	"dockside/app/util"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"os"
	"time"
)

type CloneRepoResult struct {
	Directory string
	Commit    string
}

type CloneRepoRequest struct {
	Progress sideband.Progress
	UseCache bool
}

var repoCloneCache = expirable.NewLRU[string, *CloneRepoResult](100, nil, time.Second*30)

func (bm *DockerBuildMeta) CloneRepo(request CloneRepoRequest) (*CloneRepoResult, error) {

	hash := util.HashString(
		fmt.Sprintf("%s-%s", bm.RepositoryUrl, bm.GithubAccessToken),
	)

	// pull from cache if we can
	if request.UseCache {
		if cached, ok := repoCloneCache.Get(hash); ok {
			logger.InfoWithFields("Using cached repo clone", map[string]any{
				"hash": hash,
				"repo": bm.RepositoryUrl,
			})
			return cached, nil
		}
	}

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

	result := &CloneRepoResult{
		Directory: tempDir,
		Commit:    commitHash,
	}

	if request.UseCache {
		repoCloneCache.Add(hash, result)
	}

	return result, nil
}
