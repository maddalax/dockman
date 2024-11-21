package app

import (
	"dockside/app/logger"
	"dockside/app/util/json2"
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"slices"
	"strings"
	"time"
)

type JobMetricsManager struct {
	kv *KvClient
}

type JobMetric struct {
	JobName         string
	Status          string
	LastRan         time.Time
	Interval        time.Duration
	TotalRuns       int
	LastRunDuration time.Duration
}

func NewJobMetricsManager(locator *service.Locator) *JobMetricsManager {
	return &JobMetricsManager{
		kv: KvFromLocator(locator),
	}
}

func (jb *JobMetricsManager) GetMetrics() []*JobMetric {
	metrics := make([]*JobMetric, 0)
	bucket, err := jb.kv.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "job_metrics",
	})
	if err != nil {
		logger.Error("failed to create job metrics bucket", err)
		return make([]*JobMetric, 0)
	}
	keys, err := bucket.Keys()
	if err != nil {
		return make([]*JobMetric, 0)
	}
	for _, key := range keys {
		raw, err := bucket.Get(key)
		if err != nil {
			continue
		}
		metric, err := json2.Deserialize[JobMetric](raw.Value())
		if err != nil {
			continue
		}
		metrics = append(metrics, metric)
	}

	slices.SortFunc(metrics, func(a, b *JobMetric) int {
		return strings.Compare(a.JobName, b.JobName)
	})

	return metrics
}

func (jb *JobMetricsManager) SaveJobMetric(job *Job) {
	bucket, err := jb.kv.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "job_metrics",
	})
	if err != nil {
		logger.Error("failed to create job metrics bucket", err)
		return
	}
	err = jb.kv.PutJson(bucket, job.name, &JobMetric{
		JobName:         job.name,
		Status:          job.status,
		LastRan:         job.lastRunTime,
		TotalRuns:       job.totalRuns,
		Interval:        job.interval,
		LastRunDuration: job.lastRunDuration,
	})
	if err != nil {
		logger.Error("failed to save job metric", err)
	}
}

func (jb *JobMetricsManager) OnJobStarted(job *Job) {
	logger.InfoWithFields("job started", map[string]any{
		"job_name": job.name,
	})
	jb.SaveJobMetric(job)
}

func (jb *JobMetricsManager) OnJobFinished(job *Job) {
	logger.InfoWithFields("job finished", map[string]any{
		"job_name":   job.name,
		"total_runs": job.totalRuns,
		"duration":   fmt.Sprintf("%dms", job.lastRunDuration.Milliseconds()),
	})
	jb.SaveJobMetric(job)
}

func (jb *JobMetricsManager) OnJobStopped(job *Job) {
	logger.InfoWithFields("job stopped", map[string]any{
		"job_name":   job.name,
		"total_runs": job.totalRuns,
	})
	jb.SaveJobMetric(job)
}
