package worker

import (
	"context"
	"main/job"
	"main/worker/state"
	"time"
)

type Worker interface {
	NewWorker(stateSaver *state.BasicStateSaver, maxTime time.Duration) *Worker
	ExecuteJobs(ctx context.Context) <-chan job.Result
	AppendToJobs(job job.Job)
	Stop()
}

type StateSaver interface {
	GetShutdownState(isShutdownClean bool) state.ShutdownState
	SetMemMetric(time time.Time)
	SetDiskCMetric(time time.Time)
	SetDiskDMetric(time time.Time)
	SetNetCountersMetric(time time.Time)
	SetCpuMetric(time time.Time)
}
