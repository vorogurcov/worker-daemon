package monitoring

import (
	"context"
	"fmt"
	"main/job"
	metrics2 "main/job/metrics"
	"main/worker/state"
	"time"

	cpu2 "github.com/shirou/gopsutil/cpu"
)

type CPUResult struct {
	Percent float64
}

func (r CPUResult) MetricName() string { return "cpu" }
func (r CPUResult) String() string {
	return fmt.Sprintf("cpu used: %.2f%%", r.Percent)
}

func NewCPUCallback(basicStateSaver *state.BasicStateSaver, metrics *metrics2.Metrics) job.MonitoringCallback {
	return func(ctx context.Context) (job.MonitoringResult, error) {
		percentSlice, err := cpu2.PercentWithContext(ctx, 0, false)
		if err != nil {
			return nil, err
		}
		if len(percentSlice) == 0 {
			return nil, fmt.Errorf("no cpu percent returned")
		}
		metrics.CpuUsagePercent.Set(percentSlice[0])
		basicStateSaver.SetCpuMetric(time.Now())

		return CPUResult{Percent: percentSlice[0]}, nil
	}
}
