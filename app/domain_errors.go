package app

import (
	"errors"
	"fmt"
)

var ResourceNotFoundError = errors.New("resource not found")
var ContainerExistsError = errors.New("container already exists")
var UnsupportedRunTypeError = errors.New("unsupported run type")
var BuildCancelledError = errors.New("build cancelled")
var UnknownBuildTypeError = errors.New("unknown build type")
var ResourceFailedToStopError = errors.New("resource failed to stop")
var ResourceFailedToStartError = errors.New("resource failed to start")
var DockerConnectionError = errors.New("failed to connect to docker")
var NatsKeyNotFoundError = errors.New("nats: key not found")
var ResourcePortInUseError = func(port string) error {
	return fmt.Errorf("port %s is already in use by another process, redeploy the resource to bind to a new port", port)
}
