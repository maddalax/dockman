package resources

import (
	"github.com/nats-io/nats.go"
	"paas/domain"
	"paas/kv"
)

func MapToResource(bucket nats.KeyValue) (*domain.Resource, error) {
	resource, err := kv.MustMapAllInto[domain.Resource](bucket)

	if err != nil {
		return nil, err
	}

	switch resource.RunType {
	case domain.RunTypeDockerBuild:
		resource.BuildMeta = kv.MustMapStringToStructure[domain.DockerBuildMeta](resource.BuildMeta.(string))
	case domain.RunTypeDockerRegistry:
		resource.BuildMeta = kv.MustMapStringToStructure[domain.DockerRegistryMeta](resource.BuildMeta.(string))
	default:
		resource.BuildMeta = domain.EmptyBuildMeta{}
	}

	if resource.InstancesPerServer == 0 {
		resource.InstancesPerServer = 1
	}

	return resource, nil
}
