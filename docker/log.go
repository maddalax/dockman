package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
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
		Tail:       "1000",
	})

	if err != nil {
		return err
	}

	if opts.Stdout != nil {
		go func() {
			_, err := stdcopy.StdCopy(opts.Stdout, opts.Stdout, out)
			if err != nil {
				slog.Error("error copying logs from container to stdout", slog.String("error", err.Error()))
				_ = out.Close()
			}
			fmt.Printf("done coping logs from container to stdout\n")
		}()
	}

	return nil
}
