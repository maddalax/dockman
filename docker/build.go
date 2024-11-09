package docker

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"io"
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
			cw.Writer.Write([]byte(str))
		}
	}
	return len(p), nil
}

func (c *Client) Build(out io.Writer, path string, opts types.ImageBuildOptions) error {
	ctx := context.Background()
	abs, err := filepath.Abs(path)

	buildContext, _ := archive.TarWithOptions(abs, &archive.TarOptions{})

	response, err := c.cli.ImageBuild(ctx, buildContext, opts)
	if err != nil {
		return err
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
