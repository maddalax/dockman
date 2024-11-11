package resources

import (
	"github.com/nats-io/nats.go"
	"paas/kv"
)

func MapToResource(bucket nats.KeyValue) (*Resource, error) {
	resource, err := kv.MustMapAllInto[Resource](bucket)

	if err != nil {
		return nil, err
	}

	switch resource.RunType {
	case RunTypeDockerBuild:
		resource.BuildMeta = kv.MustMapStringToStructure[DockerBuildMeta](resource.BuildMeta.(string))
	case RunTypeDockerRegistry:
		resource.BuildMeta = kv.MustMapStringToStructure[DockerRegistryMeta](resource.BuildMeta.(string))
	default:
		resource.BuildMeta = EmptyBuildMeta{}
	}

	return resource, nil
}
