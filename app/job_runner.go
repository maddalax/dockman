package app

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"sync"
	"time"
)

type Job struct {
	registry        *ServiceRegistry
	name            string
	description     string
	source          string
	locator         *service.Locator
	interval        time.Duration
	lastRunDuration time.Duration
	lastRunTime     time.Time
	status          string
	totalRuns       int
	paused          bool
	stopped         bool
	cb              func()
}

func (j *Job) Pause() {
	j.paused = true
	j.status = "paused"
	go j.registry.GetEventHandler().OnJobPaused(j)
}

func (j *Job) IsPaused() bool {
	return j.paused
}

func (j *Job) IsStopped() bool {
	return j.stopped
}

func (j *Job) Toggle() {
	if j.IsPaused() {
		j.Resume()
	} else {
		j.Pause()
	}
}

func (j *Job) Resume() {
	j.status = "running"
	j.paused = false
	go j.registry.GetEventHandler().OnJobResumed(j)
}

func (j *Job) Stop() {
	j.stopped = true
	j.status = "stopped"
	go j.registry.GetEventHandler().OnJobStopped(j)
}

type IntervalJobRunner struct {
	locator *service.Locator
	jobs    []*Job
}

func NewIntervalJobRunner(locator *service.Locator) *IntervalJobRunner {
	return &IntervalJobRunner{
		locator: locator,
	}
}

func IntervalJobRunnerFromLocator(locator *service.Locator) *IntervalJobRunner {
	return service.Get[IntervalJobRunner](locator)
}

func (jr *IntervalJobRunner) GetJob(nameAndSource string) *Job {
	for _, job := range jr.jobs {
		jobNameAndSource := fmt.Sprintf("%s-%s", job.source, job.name)
		if jobNameAndSource == nameAndSource {
			return job
		}
	}
	return nil
}

func (jr *IntervalJobRunner) Add(source string, name string, description string, duration time.Duration, job func()) {
	jr.jobs = append(jr.jobs, &Job{
		source:      source,
		name:        name,
		description: description,
		locator:     jr.locator,
		interval:    duration,
		cb:          job,
		registry:    GetServiceRegistry(jr.locator),
	})
}

func (jr *IntervalJobRunner) Start() {
	wg := sync.WaitGroup{}
	registry := GetServiceRegistry(jr.locator)
	for _, job := range jr.jobs {
		wg.Add(1)
		go func(job *Job) {
			defer wg.Done()
			for {
				if job.paused {
					time.Sleep(time.Second)
					continue
				}
				if job.stopped {
					job.status = "stopped"
					break
				}
				now := time.Now()
				job.status = "running"
				go registry.GetEventHandler().OnJobStarted(job)
				job.cb()
				job.totalRuns++
				job.status = "finished"
				go registry.GetEventHandler().OnJobFinished(job)
				job.lastRunTime = now
				job.lastRunDuration = time.Since(now)
				time.Sleep(job.interval)
			}
		}(job)
	}

	wg.Wait()

	for _, job := range jr.jobs {
		registry.GetEventHandler().OnJobStopped(job)
	}

}
