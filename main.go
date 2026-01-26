package main

import (
	"context"
	"fmt"
	"main/job"
	metrics2 "main/job/metrics"
	"main/job/monotoring"
	"main/server"
	"main/worker"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	progCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	reg := prometheus.NewRegistry()
	metrics := metrics2.NewMetrics(reg)

	go func() {
		err := server.Serve(progCtx, 8080, reg)
		if err != nil {
			fmt.Println(err)

		}
	}()

	waitingJob := job.WaitingJob{WaitTime: time.Second, WorkTime: 5 * time.Second}
	monitoringDiskCJob := job.MonitoringJob{
		Name:         "monitoringDiskCJob",
		WorkTime:     time.Hour,
		WorkInterval: time.Second,
		Callback:     monotoring.NewDiskCallback("C:", metrics),
	}
	monitoringDiskDJob := job.MonitoringJob{
		Name:         "monitoringDiskDJob",
		WorkTime:     5 * time.Second,
		WorkInterval: time.Second,
		Callback:     monotoring.NewDiskCallback("D:", metrics),
	}
	monitoringCPUJob := job.MonitoringJob{
		Name:         "monitoringCPUJob",
		WorkTime:     10 * time.Second,
		WorkInterval: 500 * time.Millisecond,
		Callback:     monotoring.NewCPUCallback(metrics),
	}
	monitoringMemJob := job.MonitoringJob{
		Name:         "monitoringMemJob",
		WorkTime:     time.Hour,
		WorkInterval: time.Second,
		Callback:     monotoring.NewMemCallback(metrics),
	}
	monitoringNetJob := job.MonitoringJob{
		Name:         "monitoringNetJob",
		WorkTime:     time.Hour,
		WorkInterval: time.Second,
		Callback:     monotoring.NewNetCallback(metrics),
	}

	jobs := []job.Job{&waitingJob, &monitoringDiskCJob, &monitoringDiskDJob,
		&monitoringCPUJob, &monitoringMemJob, &monitoringNetJob}

	basic := worker.BasicWorker{MaxWorkTime: 1 * time.Hour}
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
