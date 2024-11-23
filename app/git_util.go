package app

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
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
