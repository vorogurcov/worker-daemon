package worker

import (
	"context"
	"main/job"
	"main/worker/state"
	"time"
)

type Worker interface {
	ExecuteJobs(ctx context.Context) <-chan job.Result
	AppendToJobs(ctx context.Context, job job.Job)
	Stop() error
}

type StateSaver interface {
	GetShutdownState(isShutdownClean bool) state.ShutdownState
	SetMemMetric(time time.Time)
	SetDiskCMetric(time time.Time)
	SetDiskDMetric(time time.Time)
	SetNetCountersMetric(time time.Time)
	SetCpuMetric(time time.Time)
}
