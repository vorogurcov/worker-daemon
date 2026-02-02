package worker

import (
	"context"
	"main/job"
	"time"
)

type Worker interface {
	NewWorker(maxTime time.Duration) *Worker
	ExecuteJobs(ctx context.Context) <-chan job.Result
	AppendToJobs(job job.Job)
	Stop()
}
