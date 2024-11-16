package domain

import "errors"

var ResourceNotFoundError = errors.New("resource not found")
var ContainerExistsError = errors.New("container already exists")
var UnsupportedRunTypeError = errors.New("unsupported run type")
var BuildCancelledError = errors.New("build cancelled")
var UnknownBuildTypeError = errors.New("unknown build type")
var ResourceFailedToStopError = errors.New("resource failed to stop")
var ResourceFailedToStartError = errors.New("resource failed to start")
