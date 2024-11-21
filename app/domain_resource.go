package app

import (
	"dockside/app/util/json2"
	"encoding/json"
	"time"
)

type Resource struct {
	Id                 string            `json:"id"`
	Name               string            `json:"name"`
	Environment        string            `json:"environment"`
	RunType            RunType           `json:"run_type"`
	InstancesPerServer int               `json:"instances_per_server"`
	BuildMeta          BuildMeta         `json:"build_meta"`
	Env                map[string]string `json:"env"`
	RunStatus          RunStatus         `json:"run_status"`
	ServerDetails      []ResourceServer  `json:"server_details"`
}

type ResourceServer struct {
	ServerId   string    `json:"server_id"`
	RunStatus  RunStatus `json:"run_status"`
	LastUpdate time.Time `json:"last_update"`
}

type ResourceServerWithDetails struct {
	ResourceServer *ResourceServer
	Details        *Server
}

func (resource *Resource) MarshalJSON() ([]byte, error) {

	if resource.Env == nil {
		resource.Env = make(map[string]string)
	}

	if resource.InstancesPerServer == 0 {
		resource.InstancesPerServer = 1
	}

	if resource.ServerDetails == nil {
		resource.ServerDetails = make([]ResourceServer, 0)
	}

	buildMeta := json2.SerializeOrEmpty(resource.BuildMeta)
	serverDetails := json2.SerializeOrEmpty(resource.ServerDetails)

	return json.Marshal(map[string]interface{}{
		"id":                   resource.Id,
		"name":                 resource.Name,
		"environment":          resource.Environment,
		"run_type":             resource.RunType,
		"instances_per_server": resource.InstancesPerServer,
		"build_meta":           json.RawMessage(buildMeta),
		"env":                  resource.Env,
		"run_status":           resource.RunStatus,
		"server_details":       json.RawMessage(serverDetails),
	})
}

func (resource *Resource) UnmarshalJSON(data []byte) error {
	temp := make(map[string]interface{})
	err := json.Unmarshal(data, &temp)

	if err != nil {
		return err
	}

	resource.Id = temp["id"].(string)
	resource.Name = temp["name"].(string)
	resource.Environment = temp["environment"].(string)
	resource.InstancesPerServer = int(temp["instances_per_server"].(float64))
	resource.RunStatus = RunStatus(temp["run_status"].(float64))
	resource.RunType = RunType(temp["run_type"].(float64))

	env, ok := temp["env"].(map[string]interface{})

	if ok {
		for k, v := range env {
			resource.Env[k] = v.(string)
		}
	}

	buildMeta := temp["build_meta"].(map[string]interface{})
	switch resource.RunType {
	case RunTypeDockerBuild:
		resource.BuildMeta = &DockerBuildMeta{}
		serialized := json2.SerializeOrEmpty(buildMeta)
		err = json.Unmarshal(serialized, &resource.BuildMeta)
	case RunTypeDockerRegistry:
		resource.BuildMeta = &DockerRegistryMeta{
			Image: buildMeta["image"].(string),
		}
	default:
		resource.BuildMeta = &EmptyBuildMeta{}
	}

	serverDetails, ok := temp["server_details"].([]interface{})

	if ok {
		for _, detail := range serverDetails {
			s, err := json2.Serialize(detail)
			if err == nil {
				s2, err := json2.Deserialize[ResourceServer](s)
				if err == nil {
					resource.ServerDetails = append(resource.ServerDetails, *s2)
				}
			}
		}
	}

	return nil
}

type RunStatus = int

const (
	RunStatusUnknown RunStatus = iota
	RunStatusNotRunning
	RunStatusRunning
	RunStatusPartiallyRunning
)

func NewResource(id string) *Resource {
	resource := &Resource{
		Id: id,
	}
	return resource
}

type RunType int

const (
	RunTypeUnknown RunType = iota
	RunTypeDockerBuild
	RunTypeDockerRegistry
)

type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
