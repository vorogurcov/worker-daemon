package worker

import (
	"context"
	"main/job"
	"sync"
	"time"
)

type BasicWorker struct {
	workerJobs  chan job.Job
	MaxWorkTime time.Duration
	QueueSize   int
	wg          sync.WaitGroup
}

func NewWorker(maxTime time.Duration, queueSize int) *BasicWorker {
	return &BasicWorker{
		workerJobs:  make(chan job.Job, queueSize),
		QueueSize:   queueSize,
		MaxWorkTime: maxTime,
	}
}

func (bw *BasicWorker) ExecuteJobs(ctx context.Context) <-chan error {
	chErr := make(chan error, bw.QueueSize)

	go func() {
		workerCtx, cancel := context.WithTimeout(ctx, bw.MaxWorkTime)
		defer cancel()

		go func() {
			<-workerCtx.Done()
			bw.Stop()
		}()

		for j := range bw.workerJobs {
			bw.wg.Add(1)
			go func(jobToRun job.Job) {
				defer bw.wg.Done()

				if err := jobToRun.Do(workerCtx); err != nil {
					chErr <- err
				}
			}(j)
		}

		bw.wg.Wait()
		close(chErr)
	}()

	return chErr
}

func (bw *BasicWorker) AppendToJobs(job job.Job) {
	bw.workerJobs <- job
}

func (bw *BasicWorker) Stop() {
	close(bw.workerJobs)
}
