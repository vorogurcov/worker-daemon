package main

import (
	"context"
	"fmt"
	"main/job"
	"main/worker"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	progCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	waitingJob := job.WaitingJob{WaitTime: time.Second, WorkTime: 5 * time.Second}
	monitoringJob := job.MonitoringJob{
		WorkTime:     time.Hour,
		WorkInterval: time.Second,
	}

	jobs := []job.Job{&waitingJob, &monitoringJob}

	basic := worker.BasicWorker{MaxWorkTime: 24 * time.Hour}
	if err := basic.SetJobs(jobs); err != nil {
		fmt.Println(err)
	}

	errChan := basic.ExecuteJobs(progCtx)

	for {
		select {
		case err, ok := <-errChan:
			if !ok {
				fmt.Println("finished all jobs!")
				return
			}
			if err != nil {
				fmt.Println(err)
			}

		}
	}

}
