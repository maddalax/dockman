package docker

import (
	"context"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func Connect() (*Client, error) {
	env := client.FromEnv
	cli, err := client.NewClientWithOpts(env,
		// TODO
		client.WithHost("unix:///Users/maddox/.docker/run/docker.sock"),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	_, err = cli.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return &Client{
		cli: cli,
	}, nil
}

//func main() {
//	env := client.FromEnv
//
//	cli, err := client.NewClientWithOpts(env,
//		client.WithHost("unix:///Users/maddox/.docker/run/docker.sock"),
//		client.WithAPIVersionNegotiation(),
//	)
//	if err != nil {
//		panic(err)
//	}
//
//	ctx := context.Background()
//
//	defer cli.Close()
//
//	cli.ImageBuild(ctx, os.Stdin, types.ImageBuildOptions{})
//
//	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", image.PullOptions{})
//	if err != nil {
//		panic(err)
//	}
//	io.Copy(os.Stdout, reader)
//
//	resp, err := cli.ContainerCreate(ctx, &container.Config{
//		Image: "alpine",
//		Cmd:   []string{"echo", "hello world"},
//	}, nil, nil, nil, "")
//	if err != nil {
//		panic(err)
//	}
//
//	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
//		panic(err)
//	}
//
//	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
//	select {
//	case err := <-errCh:
//		if err != nil {
//			panic(err)
//		}
//	case <-statusCh:
//	}
//
//	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
//	if err != nil {
//		panic(err)
//	}
//
//	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
//}
