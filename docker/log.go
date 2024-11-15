package docker

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"log/slog"
	"strconv"
	"time"
)

type StreamLogsOptions struct {
	Stdout io.WriteCloser
	Since  time.Time
}

func (c *Client) StreamLogs(containerId string, opts StreamLogsOptions) error {
	ctx := context.Background()

	logOpts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "1000",
		Timestamps: true,
	}

	if !opts.Since.IsZero() {
		logOpts.Since = strconv.FormatInt(opts.Since.Unix(), 10)
	}

	out, err := c.cli.ContainerLogs(ctx, containerId, logOpts)

	if err != nil {
		return err
	}

	if opts.Stdout != nil {
		_, err := stdcopy.StdCopy(opts.Stdout, opts.Stdout, out)
		if err != nil {
			slog.Error("failed to copy logs from docker container", slog.String("containerId", containerId), slog.String("error", err.Error()))
		}
		_ = out.Close()
		_ = opts.Stdout.Close()
	}

	return nil
}
