package monitoring

import (
	"context"
	"fmt"
	"main/job"
	metrics2 "main/job/metrics"
	"main/worker/state"
	"time"

	"github.com/shirou/gopsutil/disk"
)

type DiskResult struct {
	Path        string
	UsedPercent float64
}

func (r DiskResult) MetricName() string { return "disk" }
func (r DiskResult) String() string {
	return fmt.Sprintf("disk(%s) used: %.2f%%", r.Path, r.UsedPercent)
}

func NewDiskCallback(basicStateSaver *state.BasicStateSaver, path string, metrics *metrics2.Metrics) job.MonitoringCallback {
	return func(ctx context.Context) (job.MonitoringResult, error) {
		u, err := disk.UsageWithContext(ctx, path)
		if err != nil {
			return nil, err
		}
		if path == "C:" {
			metrics.DiskCUsagePercent.Set(u.UsedPercent)
			basicStateSaver.SetDiskCMetric(time.Now())

		} else if path == "D:" {
			metrics.DiskDUsagePercent.Set(u.UsedPercent)
			basicStateSaver.SetDiskDMetric(time.Now())

		}

		return DiskResult{Path: path, UsedPercent: u.UsedPercent}, nil
	}
}
