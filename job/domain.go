package job

import "context"

type Job interface {
	Do(ctx context.Context) error
}

//type MonitoringStatisticsCallback func(any interface{}) string
