package worker

import (
	"context"
	"main/job"
)

type Worker interface {
	ExecuteJobs(ctx context.Context) <-chan error
	SetJobs(ctx context.Context, jobs []job.Job)
}
