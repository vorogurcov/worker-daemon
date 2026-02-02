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

func (bw *BasicWorker) ExecuteJobs(ctx context.Context) <-chan job.Result {
	resCh := make(chan job.Result, 10*bw.QueueSize)

	go func() {
		workerCtx, cancel := context.WithTimeout(ctx, bw.MaxWorkTime)
		defer cancel()
		defer close(resCh)
		
		go func() {
			<-workerCtx.Done()
			bw.Stop()
		}()

		for j := range bw.workerJobs {
			bw.wg.Add(1)
			go func(jobToRun job.Job) {
				defer bw.wg.Done()
				for res := range jobToRun.Do(workerCtx) {
					resCh <- res
				}

			}(j)
		}

		bw.wg.Wait()

	}()

	return resCh
}

func (bw *BasicWorker) AppendToJobs(job job.Job) {
	bw.workerJobs <- job
}

func (bw *BasicWorker) Stop() {
	close(bw.workerJobs)
}
