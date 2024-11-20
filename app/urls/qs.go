package urls

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
)

func WithQs(url string, pairs ...string) string {
	if len(pairs) == 0 {
		return url
	}
	qs := h.NewQs(pairs...).ToString()
	return fmt.Sprintf("%s?%s", url, qs)
}

func ResourceUrl(id string) string {
	return WithQs("/resource", "id", id)
}

func ResourceServersUrl(id string) string {
	return WithQs("/resource/servers", "id", id)
}

func ServerUrl(id string) string {
	return WithQs("/server", "id", id)
}

func ResourceStartDeploymentPath(resourceId string, buildId string) string {
	return WithQs("/resource/deployment/new", "resourceId", resourceId, "buildId", buildId)
}

func ResourceDeploymentLogUrl(id string, buildId string) string {
	return WithQs("/resource/deployment/build-log", "id", id, "buildId", buildId)
}

func ResourceRunLogUrl(id string) string {
	return WithQs("/resource/deployment/run-log", "id", id)
}

func ResourceEnvironmentUrl(id string) string {
	return WithQs("/resource/deployment/environment", "id", id)
}

func ResourceDeploymentUrl(id string) string {
	return WithQs("/resource/deployment", "id", id)
}

func NewResourceUrl() string {
	return WithQs("/resource/create")
}

func NewServerUrl() string {
	return WithQs("/server/create")
}
