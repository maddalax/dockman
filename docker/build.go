package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"io"
	"path/filepath"
)

func (c *Client) Build(out io.Writer, path string, opts types.ImageBuildOptions) error {
	ctx := context.Background()
	abs, err := filepath.Abs(path)

	buildContext, _ := archive.TarWithOptions(abs, &archive.TarOptions{})

	response, err := c.cli.ImageBuild(ctx, buildContext, opts)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Read the output of the build process
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	return nil
}
