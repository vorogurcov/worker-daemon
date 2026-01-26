package monotoring

import (
	"context"
	"fmt"
	"main/job"
	metrics2 "main/job/metrics"
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

func NewDiskCallback(path string, metrics *metrics2.Metrics) job.MonitoringCallback {
	return func(ctx context.Context) (job.MonitoringResult, error) {
		cctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		u, err := disk.UsageWithContext(cctx, path)
		if err != nil {
			return nil, err
		}
		if path == "C:" {
			metrics.DiskCUsagePercent.Set(u.UsedPercent)
		} else if path == "D:" {
			metrics.DiskDUsagePercent.Set(u.UsedPercent)
		}

		return DiskResult{Path: path, UsedPercent: u.UsedPercent}, nil
	}
}
