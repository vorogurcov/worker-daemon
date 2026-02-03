package state

import "time"

type BasicStateSaver struct {
	SchemaVer      int
	cpuTime        time.Time
	memTime        time.Time
	diskCTime      time.Time
	diskDTime      time.Time
	netCounterTime time.Time
}

func (b *BasicStateSaver) GetShutdownState(isShutdownClean bool) ShutdownState {

	return ShutdownState{
		SchemaVersion:     b.SchemaVer,
		LastShutdownClean: isShutdownClean,
		ShutdownTimeUTC:   time.Now(),
		LastCollect: LastCollectTimeUTCByMetric{
			Cpu:        b.cpuTime,
			Mem:        b.memTime,
			DiskC:      b.diskCTime,
			DiskD:      b.diskDTime,
			NetCounter: b.netCounterTime,
		},
	}
}

func (b *BasicStateSaver) SetMemMetric(time time.Time) {
	b.memTime = time
}

func (b *BasicStateSaver) SetDiskCMetric(time time.Time) {
	b.diskCTime = time
}

func (b *BasicStateSaver) SetDiskDMetric(time time.Time) {
	b.diskDTime = time

}

func (b *BasicStateSaver) SetNetCountersMetric(time time.Time) {
	b.netCounterTime = time
}

func (b *BasicStateSaver) SetCpuMetric(time time.Time) {
	b.cpuTime = time
}
