package app

func (resource *Resource) RunTypeFormatted() string {
	switch resource.RunType {
	case RunTypeDockerRegistry:
		return "Docker Registry"
	case RunTypeDockerBuild:
		return "Dockerfile build"
	default:
		return "Unknown"
	}
}
