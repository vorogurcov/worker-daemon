package monitoring

import (
	"context"
	"fmt"
	"main/job"
	metrics2 "main/job/metrics"
	"main/worker/state"
	"time"

	"github.com/shirou/gopsutil/net"
)

type NetResult struct {
	BytesSent uint64
	BytesRecv uint64
}

func (r NetResult) MetricName() string { return "net" }
func (r NetResult) String() string {
	return fmt.Sprintf("total bytes sent: %v MiB, recv: %v MiB",
		r.BytesSent/1024/1024, r.BytesRecv/1024/1024)
}

func NewNetCallback(basicStateSaver *state.BasicStateSaver, metrics *metrics2.Metrics) job.MonitoringCallback {
	return func(ctx context.Context) (job.MonitoringResult, error) {
		counters, err := net.IOCountersWithContext(ctx, false)
		if err != nil {
			return nil, err
		}
		if len(counters) == 0 {
			return nil, fmt.Errorf("no net counters")
		}
		metrics.NetCounter.Set(float64(counters[0].BytesSent / 1024 / 1024))
		basicStateSaver.SetNetCountersMetric(time.Now())

		return NetResult{
			BytesSent: counters[0].BytesSent,
			BytesRecv: counters[0].BytesRecv,
		}, nil
	}
}
