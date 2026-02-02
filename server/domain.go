package server

import "time"

type CreateJobDto struct {
	Type         string        `json:"type"` // "WaitingJob" | "MonitoringCPUJob"
	Name         string        `json:"name"`
	WorkTime     time.Duration `json:"work_time"`
	WorkInterval time.Duration `json:"work_interval"`
}
