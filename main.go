package main

import (
	"context"
	"fmt"
	"main/job"
	metrics2 "main/job/metrics"
	"main/job/monitoring"
	"main/server"
	"main/worker"
	"main/worker/state"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	progCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	reg := prometheus.NewRegistry()
	metrics := metrics2.NewMetrics(reg)

	b := state.BasicStateSaver{SchemaVer: 1}

	maxWorkTime := time.Duration(0)

	waitingJob := job.WaitingJob{WorkInterval: time.Second, WorkTime: 5 * time.Second}
	monitoringDiskCJob := job.MonitoringJob{
		Name:         "monitoringDiskCJob",
		WorkTime:     30 * time.Second,
		WorkInterval: time.Second,
		Callback:     monitoring.NewDiskCallback(&b, "C:", metrics),
	}
	monitoringDiskDJob := job.MonitoringJob{
		Name:         "monitoringDiskDJob",
		WorkTime:     30 * time.Second,
		WorkInterval: time.Second,
		Callback:     monitoring.NewDiskCallback(&b, "D:", metrics),
	}
	monitoringCPUJob := job.MonitoringJob{
		Name:         "monitoringCPUJob",
		WorkTime:     30 * time.Second,
		WorkInterval: 500 * time.Millisecond,
		Callback:     monitoring.NewCPUCallback(&b, metrics),
	}
	monitoringMemJob := job.MonitoringJob{
		Name:         "monitoringMemJob",
		WorkTime:     30 * time.Second,
		WorkInterval: time.Second,
		Callback:     monitoring.NewMemCallback(&b, metrics),
	}
	monitoringNetJob := job.MonitoringJob{
		Name:         "monitoringNetJob",
		WorkTime:     30 * time.Second,
		WorkInterval: time.Second,
		Callback:     monitoring.NewNetCallback(&b, metrics),
	}

	jobs := []job.Job{&waitingJob, &monitoringDiskCJob, &monitoringDiskDJob,
		&monitoringCPUJob, &monitoringMemJob, &monitoringNetJob}

	basic := worker.NewWorker(&b, maxWorkTime, 100)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.Serve(progCtx, maxWorkTime, 8080, reg, metrics, &b, basic)
		if err != nil {
			fmt.Println(err)
		}
	}()

	resCh := basic.ExecuteJobs(progCtx)

	for _, j := range jobs {
		basic.AppendToJobs(progCtx, j)
	}

	go func() {
		time.Sleep(8 * time.Second)
		waitJob := job.WaitingJob{WorkInterval: time.Second, WorkTime: 5 * time.Second}
		basic.AppendToJobs(progCtx, &waitJob)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for res := range resCh {
			if res.Error != nil {
				fmt.Printf("%v", res.Error)
			}
			if res.Value != nil {
				fmt.Printf("%v", res.Value)
			}
		}
		fmt.Println("finished all jobs!")
	}()

	wg.Wait()
}
