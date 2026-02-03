package state

import (
	"time"
)

type LastCollectTimeUTCByMetric struct {
	Cpu        time.Time `json:"cpu"`
	Mem        time.Time `json:"mem"`
	DiskC      time.Time `json:"disk_c"`
	DiskD      time.Time `json:"disk_d"`
	NetCounter time.Time `json:"net_counter"`
}

type ShutdownState struct {
	SchemaVersion     int                        `json:"schema_version"`
	LastShutdownClean bool                       `json:"last_shutdown_clean"`
	ShutdownTimeUTC   time.Time                  `json:"shutdown_time_utc"`
	LastCollect       LastCollectTimeUTCByMetric `json:"last_collect"`
}
