package monitoring

import (
	"context"
	"fmt"
	"main/job"
	metrics2 "main/job/metrics"
	"main/worker/state"
	"time"

	"github.com/shirou/gopsutil/mem"
)

type MemResult struct {
	UsedPercent float64
}

func (r MemResult) MetricName() string { return "memory" }
func (r MemResult) String() string {
	return fmt.Sprintf("memory used: %.2f%%", r.UsedPercent)
}

func NewMemCallback(basicStateSaver *state.BasicStateSaver, metrics *metrics2.Metrics) job.MonitoringCallback {
	return func(ctx context.Context) (job.MonitoringResult, error) {
		v, err := mem.VirtualMemoryWithContext(ctx)
		if err != nil {
			return nil, err
		}
		metrics.MemUsagePercent.Set(v.UsedPercent)
		basicStateSaver.SetMemMetric(time.Now())

		return MemResult{UsedPercent: v.UsedPercent}, nil
	}
}
