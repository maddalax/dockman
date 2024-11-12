package subject

import "fmt"

type Subject = string

func BuildLogForResource(id string, buildId string) string {
	return fmt.Sprintf("build.log-%s-%s", id, buildId)
}

var ResourceCreated = "resource.created"
