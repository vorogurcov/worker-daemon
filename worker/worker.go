package worker

import (
	"context"
	"errors"
	"fmt"
	"main/job"
	"sync"
	"time"
)

type BasicWorker struct {
	workerJobs  []job.Job
	MaxWorkTime time.Duration
}

func (bw *BasicWorker) ExecuteJobs(ctx context.Context) <-chan error {
	workerCtx, cancel := context.WithTimeout(ctx, bw.MaxWorkTime)

	chErr := make(chan error, len(bw.workerJobs))

	if bw.workerJobs == nil {
		chErr <- errors.New("workerJobs are not set, call SetJobs to set it")
		return chErr
	}

	var wg sync.WaitGroup

	wg.Add(len(bw.workerJobs))

	for _, workerJob := range bw.workerJobs {

		go func(j job.Job) {
			defer wg.Done()
			if err := j.Do(workerCtx); err != nil {
				chErr <- fmt.Errorf("error in workerJob: %v", err)
			}
			chErr <- nil
		}(workerJob)
	}
	go func() {
		wg.Wait()
		cancel()
		close(chErr)
	}()

	return chErr
}

func (bw *BasicWorker) SetJobs(jobs []job.Job) error {
	bw.workerJobs = jobs
	return nil
}
