package app

import (
	"dockman/app/logger"
	"dockman/app/util/json2"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"slices"
	"strings"
	"time"
)

type JobMetricsManager struct {
	locator *service.Locator
}

type JobMetric struct {
	JobName         string
	JobSource       string
	JobDescription  string
	Status          string
	JobPaused       bool
	LastRan         time.Time
	Interval        time.Duration
	TotalRuns       int
	LastRunDuration time.Duration
}

func NewJobMetricsManager(locator *service.Locator) *JobMetricsManager {
	return &JobMetricsManager{
		locator: locator,
	}
}

func (jb *JobMetricsManager) GetMetrics() []*JobMetric {
	metrics := make([]*JobMetric, 0)
	registry := GetServiceRegistry(jb.locator)
	kv := registry.KvClient()
	bucket, err := kv.GetOrCreateBucket(&nats.KeyValueConfig{
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
		if metric.JobSource == "" {
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
	registry := GetServiceRegistry(jb.locator)
	kv := registry.KvClient()
	bucket, err := kv.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "job_metrics",
	})
	if err != nil {
		logger.Error("failed to create job metrics bucket", err)
		return
	}
	err = kv.PutJson(bucket, job.name, &JobMetric{
		JobName:         job.name,
		JobSource:       job.source,
		JobDescription:  job.description,
		Status:          job.status,
		LastRan:         job.lastRunTime,
		TotalRuns:       job.totalRuns,
		Interval:        job.interval,
		JobPaused:       job.IsPaused(),
		LastRunDuration: job.lastRunDuration,
	})
	if err != nil {
		logger.Error("failed to save job metric", err)
	}
}

func (jb *JobMetricsManager) OnJobStarted(job *Job) {
	jb.SaveJobMetric(job)
}

func (jb *JobMetricsManager) OnJobFinished(job *Job) {
	jb.SaveJobMetric(job)
}

func (jb *JobMetricsManager) OnJobStopped(job *Job) {
	jb.SaveJobMetric(job)
}

func (jb *JobMetricsManager) OnJobPaused(job *Job) {
	jb.SaveJobMetric(job)
}

func (jb *JobMetricsManager) OnJobResumed(job *Job) {
	jb.SaveJobMetric(job)
}
