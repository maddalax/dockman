package subject

import "fmt"

type Subject = string

func BuildLogForResource(id string, buildId string) string {
	return fmt.Sprintf("build.log-%s-%s", id, buildId)
}

func RunLogsForResource(id string) string {
	return fmt.Sprintf("run.log-%s", id)
}

var ResourceCreated = "resource.created"
