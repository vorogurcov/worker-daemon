package monotoring

import (
	"context"
	"fmt"
	"main/job"
	metrics2 "main/job/metrics"
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

func NewMemCallback(metrics *metrics2.Metrics) job.MonitoringCallback {
	return func(ctx context.Context) (job.MonitoringResult, error) {
		cctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		v, err := mem.VirtualMemoryWithContext(cctx)
		if err != nil {
			return nil, err
		}
		metrics.MemUsagePercent.Set(v.UsedPercent)
		return MemResult{UsedPercent: v.UsedPercent}, nil
	}
}
