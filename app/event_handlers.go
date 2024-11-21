package app

import (
	"dockside/app/logger"
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
)

type EventHandler struct {
	locator           *service.Locator
	jobMetricsManager *JobMetricsManager
}

func NewEventHandler(locator *service.Locator) *EventHandler {
	return &EventHandler{
		locator:           locator,
		jobMetricsManager: NewJobMetricsManager(locator),
	}
}

func (eh *EventHandler) OnJobStarted(job *Job) {
	logger.InfoWithFields("job started", map[string]any{
		"job_name": job.name,
	})
	eh.jobMetricsManager.OnJobStarted(job)
}

func (eh *EventHandler) OnJobFinished(job *Job) {
	logger.InfoWithFields("job finished", map[string]any{
		"job_name":   job.name,
		"total_runs": job.totalRuns,
		"duration":   fmt.Sprintf("%dms", job.lastRunDuration.Milliseconds()),
	})
	eh.jobMetricsManager.OnJobFinished(job)
}

func (eh *EventHandler) OnJobStopped(job *Job) {
	logger.InfoWithFields("job stopped", map[string]any{
		"job_name":   job.name,
		"total_runs": job.totalRuns,
	})
	eh.jobMetricsManager.OnJobStopped(job)
}

func (eh *EventHandler) OnServerDisconnected(server *Server) {
	logger.InfoWithFields("server disconnected", map[string]any{
		"server_id": server.Id,
		"name":      server.FormattedName(),
	})
}

func (eh *EventHandler) OnServerConnected(server *Server) {
	logger.InfoWithFields("server connected", map[string]any{
		"server_id": server.Id,
		"name":      server.FormattedName(),
	})
}

func (eh *EventHandler) OnServerDetached(serverId string, resource *Resource) {
	logger.InfoWithFields("server detached from resource", map[string]any{
		"server_id":     serverId,
		"resource_id":   resource.Id,
		"resource_name": resource.Name,
	})
}

func (eh *EventHandler) OnResourceStatusChange(resource *Resource, status RunStatus) {
	logger.InfoWithFields("resource status changed", map[string]any{
		"resource_id":   resource.Id,
		"resource_name": resource.Name,
		"new_status":    status,
	})
}
