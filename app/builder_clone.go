package app

import (
	"dockside/app/logger"
	"dockside/app/util"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/maddalax/htmgo/framework/h"
	"os"
	"strings"
	"time"
)

type CloneRepoResult struct {
	Directory string
	Commit    string
	Repo      *git.Repository
}

type CloneRepoRequest struct {
	Progress     sideband.Progress
	UseCache     bool
	SingleBranch bool
	BranchName   string
}

var repoCloneCache = expirable.NewLRU[string, *CloneRepoResult](100, nil, time.Second*30)

func (bm *DockerBuildMeta) CloneRepo(request CloneRepoRequest) (*CloneRepoResult, error) {

	hash := util.HashString(
		fmt.Sprintf("%s-%s-%s-%v", bm.RepositoryUrl, bm.GithubAccessToken, request.BranchName, request.SingleBranch),
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
		URL:           bm.RepositoryUrl,
		Auth:          opts.Auth,
		Progress:      request.Progress,
		Depth:         1,
		RemoteName:    "origin",
		ReferenceName: plumbing.ReferenceName(request.BranchName),
		SingleBranch:  request.SingleBranch,
	})

	if err != nil {
		if errors.Is(err, git.NoMatchingRefSpecError{}) {
			return nil, errors.New(fmt.Sprintf("branch '%s' not found in repository", bm.DeploymentBranch))
		}
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
		Repo:      repo,
	}

	if request.UseCache {
		repoCloneCache.Add(hash, result)
	}

	return result, nil
}

func (bm *DockerBuildMeta) ListRemoteBranches() ([]string, error) {
	remote := git.NewRemote(nil, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{bm.RepositoryUrl},
	})

	opts := &git.ListOptions{}

	if bm.GithubAccessToken != "" {
		opts.Auth = &http.BasicAuth{
			Username: "dockside",
			Password: bm.GithubAccessToken,
		}
	}

	refs, err := remote.List(opts)

	if err != nil {
		return nil, err
	}

	names := make([]string, 0)

	for _, ref := range refs {
		if ref.Name().IsRemote() || ref.Name().IsBranch() {
			name := ref.Name().Short()
			name = strings.TrimPrefix(name, "origin/")
			names = append(names, name)
		}
	}

	return h.Unique(names, func(item string) string {
		return item
	}), nil
}

func (bm *DockerBuildMeta) GetLatestCommitOnRemote() (string, error) {
	remote := git.NewRemote(nil, &config.RemoteConfig{
		URLs: []string{bm.RepositoryUrl},
	})

	opts := &git.ListOptions{}

	if bm.GithubAccessToken != "" {
		opts.Auth = &http.BasicAuth{
			Username: "dockside",
			Password: bm.GithubAccessToken,
		}
	}

	refs, err := remote.List(opts)

	if err != nil {
		return "", err
	}

	for _, ref := range refs {
		if ref.Name() == plumbing.NewBranchReferenceName(bm.DeploymentBranch) {
			return ref.Hash().String(), nil
		}
	}

	return "", errors.New("branch not found")
}
