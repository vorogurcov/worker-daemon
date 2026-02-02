package main

import (
	"context"
	"fmt"
	"main/job"
	metrics2 "main/job/metrics"
	"main/job/monitoring"
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

	waitingJob := job.WaitingJob{WorkInterval: time.Second, WorkTime: 5 * time.Second}
	monitoringDiskCJob := job.MonitoringJob{
		Name:         "monitoringDiskCJob",
		WorkTime:     time.Hour,
		WorkInterval: time.Second,
		Callback:     monitoring.NewDiskCallback("C:", metrics),
	}
	monitoringDiskDJob := job.MonitoringJob{
		Name:         "monitoringDiskDJob",
		WorkTime:     5 * time.Second,
		WorkInterval: time.Second,
		Callback:     monitoring.NewDiskCallback("D:", metrics),
	}
	monitoringCPUJob := job.MonitoringJob{
		Name:         "monitoringCPUJob",
		WorkTime:     10 * time.Second,
		WorkInterval: 500 * time.Millisecond,
		Callback:     monitoring.NewCPUCallback(metrics),
	}
	monitoringMemJob := job.MonitoringJob{
		Name:         "monitoringMemJob",
		WorkTime:     time.Hour,
		WorkInterval: time.Second,
		Callback:     monitoring.NewMemCallback(metrics),
	}
	monitoringNetJob := job.MonitoringJob{
		Name:         "monitoringNetJob",
		WorkTime:     time.Hour,
		WorkInterval: time.Second,
		Callback:     monitoring.NewNetCallback(metrics),
	}

	jobs := []job.Job{&waitingJob, &monitoringDiskCJob, &monitoringDiskDJob,
		&monitoringCPUJob, &monitoringMemJob, &monitoringNetJob}

	basic := worker.NewWorker(1*time.Hour, 100)

	go func() {
		err := server.Serve(progCtx, 8080, reg, metrics, basic)
		if err != nil {
			fmt.Println(err)

		}
	}()

	errChan := basic.ExecuteJobs(progCtx)

	for _, j := range jobs {
		basic.AppendToJobs(j)
	}

	go func() {
		time.Sleep(8 * time.Second)
		waitJob := job.WaitingJob{WorkInterval: time.Second, WorkTime: 5 * time.Second}
		basic.AppendToJobs(&waitJob)
	}()

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
