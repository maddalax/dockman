package app

import (
	"dockman/app/util/fileio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ParsedDockerFile struct {
	FromPort int
}

func (bm *DockerBuildMeta) SetDefaultsFromRepository() {
	if bm.RepositoryUrl == "" {
		return
	}

	clone, err := bm.CloneRepo(CloneRepoRequest{
		UseCache:     true,
		Progress:     os.Stdout,
		SingleBranch: false,
	})

	if err != nil {
		return
	}

	parsedDockerFile, err := bm.parseDockerFile(clone.Directory)

	if err == nil && parsedDockerFile.FromPort > 0 {
		bm.ExposedPort = parsedDockerFile.FromPort
	}

	// try to find the default branch for deployment
	head, err := clone.Repo.Head()
	if err == nil && head.Name().IsBranch() {
		bm.DeploymentBranch = head.Name().Short()
	} else {
		branches, err := bm.ListRemoteBranches()
		if err == nil {
			for _, branch := range branches {
				if branch == "master" || branch == "main" {
					bm.DeploymentBranch = branch
				}
			}
		}
	}
}

func (bm *DockerBuildMeta) parseDockerFile(repoDir string) (*ParsedDockerFile, error) {
	// ensure its a valid dockerfile
	validator := ValidDockerFileValidator{
		RepositoryDir: repoDir,
		Dockerfile:    bm.Dockerfile,
	}

	err := validator.Validate()

	if err != nil {
		return nil, err
	}

	dockerFilePath, err := filepath.Abs(filepath.Join(repoDir, bm.Dockerfile))

	if err != nil {
		return nil, err
	}

	var parsed = &ParsedDockerFile{}

	err = fileio.ReadLines(dockerFilePath, func(line string) {
		lower := strings.TrimSpace(strings.ToLower(line))
		if strings.HasPrefix(lower, "expose") {
			port := strings.Split(lower, " ")
			if len(port) > 1 {
				conv, err := strconv.Atoi(port[1])
				if err == nil {
					parsed.FromPort = conv
				}
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return parsed, nil
}
