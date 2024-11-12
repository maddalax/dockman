package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"io"
	"log/slog"
)

type StreamLogsOptions struct {
	Stdout io.Writer
}

func (c *Client) StreamLogs(containerId string, opts StreamLogsOptions) error {
	ctx := context.Background()
	out, err := c.cli.ContainerLogs(ctx, containerId, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})

	if err != nil {
		return err
	}

	if opts.Stdout != nil {
		go func() {
			fmt.Printf("coping logs from container to stdout\n")
			_, err := io.Copy(opts.Stdout, out)
			if err != nil {
				slog.Error("error copying logs from container to stdout", slog.String("error", err.Error()))
				_ = out.Close()
			}
			fmt.Printf("done coping logs from container to stdout\n")
		}()
	}

	return nil
}
