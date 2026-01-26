package job

import "context"

type Job interface {
	Do(ctx context.Context) error
}

type MonitoringResult interface {
	MetricName() string
	String() string
}
type MonitoringCallback func(ctx context.Context) (MonitoringResult, error)
