package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/buildkite/terminal-to-html/v3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type CustomWriter struct {
	Writer io.Writer
}

func (cw *CustomWriter) Write(p []byte) (n int, err error) {
	str := string(p)
	str = strings.ReplaceAll(str, "\r", "\n")
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		if line == "" || line == "\n" {
			continue
		}

		if !strings.HasPrefix(line, "{") {
			cw.Writer.Write([]byte(line))
			continue
		}

		d := make(map[string]any)
		err := json.Unmarshal([]byte(line), &d)

		// not json i suppose
		if err != nil {
			cw.Writer.Write([]byte(line))
			continue
		}

		// extract stream from {"stream":"my text"}
		if v, ok := d["stream"]; ok {
			str := v.(string)
			if str[0] == '"' {
				str = str[1:]
			}
			if str[len(str)-1] == '"' {
				str = str[:len(str)-1]
			}
			if str == "\n" {
				continue
			}
			cw.Writer.Write([]byte(terminal.Render([]byte(str))))
		}
	}
	return len(p), nil
}

type BuildResponse struct {
	CancelChan chan func() error
}

func (c *DockerClient) Build(out io.Writer, path string, opts types.ImageBuildOptions, cb *BuildResponse) error {
	ctx := context.Background()
	abs, err := filepath.Abs(path)
	projectDir := filepath.Dir(filepath.Join(abs, opts.Dockerfile))

	opts.Dockerfile = filepath.Base(opts.Dockerfile)

	var ignored []string
	dockerIgnore, err := os.Open(filepath.Join(projectDir, ".dockerignore"))

	if err == nil {
		scanner := bufio.NewScanner(dockerIgnore)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			ignored = append(ignored, scanner.Text())
		}
	}

	buildContext, _ := archive.TarWithOptions(projectDir, &archive.TarOptions{
		ExcludePatterns: ignored,
	})
	defer buildContext.Close()

	response, err := c.cli.ImageBuild(ctx, buildContext, opts)

	if err != nil {
		return err
	}

	if cb != nil {
		cb.CancelChan <- func() error {
			// send a request to cancel
			_ = c.cli.BuildCancel(ctx, opts.BuildID)
			// sometimes that doesn't work, so kill the response, which forces it to end
			_ = response.Body.Close()
			return nil
		}
	}

	defer response.Body.Close()

	customWriter := CustomWriter{out}
	// Read the output of the build process
	_, err = io.Copy(&customWriter, response.Body)
	if err != nil {
		return err
	}

	return nil
}
