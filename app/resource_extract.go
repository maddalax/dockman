package app

import (
	"dockside/app/util/fileio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ParsedDockerFile struct {
	FromPort int
}

func (bm *DockerBuildMeta) ParseDockerFile() (*ParsedDockerFile, error) {
	clone, err := bm.CloneRepo(CloneRepoRequest{
		UseCache: true,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	// ensure its a valid dockerfile
	validator := ValidDockerFileValidator{
		RepositoryDir: clone.Directory,
		Dockerfile:    bm.Dockerfile,
	}

	err = validator.Validate()

	if err != nil {
		return nil, err
	}

	dockerFilePath, err := filepath.Abs(filepath.Join(clone.Directory, bm.Dockerfile))

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
