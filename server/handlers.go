package server

import (
	"context"
	"errors"
	"main/job"
	"main/job/metrics"
	"main/job/monitoring"
	"main/worker"
	"main/worker/state"
)

func CreateJob(ctx context.Context, b *state.BasicStateSaver, metrics *metrics.Metrics, w *worker.BasicWorker, createJobDto CreateJobDto) error {
	var j job.Job

	if createJobDto.Type == "WaitingJob" {
		j = &job.WaitingJob{WorkInterval: createJobDto.WorkInterval, WorkTime: createJobDto.WorkTime}
	} else if createJobDto.Type == "MonitoringCPUJob" {
		j = &job.MonitoringJob{
			Name:         "monitoringCPUJob",
			WorkTime:     createJobDto.WorkTime,
			WorkInterval: createJobDto.WorkInterval,
			Callback:     monitoring.NewCPUCallback(b, metrics),
		}
	} else {
		return errors.New("type is not supported")
	}

	w.AppendToJobs(ctx, j)
	return nil
}
