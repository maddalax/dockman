package internal

import "time"

type DeploymentStatus string

const (
	DeploymentStatusPending   DeploymentStatus = "pending"
	DeploymentStatusRunning   DeploymentStatus = "running"
	DeploymentStatusSucceeded DeploymentStatus = "succeeded"
	DeploymentStatusFailed    DeploymentStatus = "failed"
)

type CreateDeploymentRequest struct {
	ResourceId string
	BuildId    string
}

type UpdateDeploymentStatusRequest struct {
	ResourceId string
	BuildId    string
	Status     DeploymentStatus
}

type Deployment struct {
	ResourceId string           `json:"resourceId"`
	CreatedAt  time.Time        `json:"createdAt"`
	BuildId    string           `json:"buildId"`
	Status     DeploymentStatus `json:"status"`
}
