package job

import "context"

type Result struct {
	Value any
	Error error
}

type Job interface {
	Do(ctx context.Context) <-chan Result
}

type MonitoringResult interface {
	MetricName() string
	String() string
}
type MonitoringCallback func(ctx context.Context) (MonitoringResult, error)
