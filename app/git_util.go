package app

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/maddalax/htmgo/framework/h"
	"os"
	"strings"
)

func GetLatestCommitOnRemote(repoUrl, branchName string) (string, error) {
	remote := git.NewRemote(nil, &config.RemoteConfig{
		URLs: []string{repoUrl},
	})

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}

	for _, ref := range refs {
		if ref.Name() == plumbing.NewBranchReferenceName(branchName) {
			return ref.Hash().String(), nil
		}
	}

	return "", errors.New("branch not found")
}

func ListRemoteBranches(repo *git.Repository) ([]string, error) {
	err := repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})

	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil, err
	}

	// Get the branch references
	refs, err := repo.References()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsRemote() || ref.Name().IsBranch() {
			name := ref.Name().Short()
			name = strings.TrimPrefix(name, "origin/")
			names = append(names, name)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return h.Unique(names, func(item string) string {
		return item
	}), nil
}
