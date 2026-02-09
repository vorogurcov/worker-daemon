package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type BasicStateSaver struct {
	mu             sync.Mutex
	SchemaVer      int
	cpuTime        time.Time
	memTime        time.Time
	diskCTime      time.Time
	diskDTime      time.Time
	netCounterTime time.Time
}

func (b *BasicStateSaver) GetLastShutdownState() (ShutdownState, error) {
	const dir = "saves"
	lastPath := filepath.Join(dir, "last-save.json")

	var state ShutdownState

	f, err := os.Open(lastPath)
	if err != nil {
		return state, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&state); err != nil {
		return state, err
	}

	return state, nil
}

func (b *BasicStateSaver) GetShutdownState(isShutdownClean bool) ShutdownState {
	prevState, err := b.GetLastShutdownState()
	b.mu.Lock()
	defer b.mu.Unlock()
	if err == nil {
		if b.cpuTime.IsZero() {
			b.cpuTime = prevState.LastCollect.Cpu
		}
		if b.memTime.IsZero() {
			b.memTime = prevState.LastCollect.Mem
		}
		if b.diskCTime.IsZero() {
			b.diskCTime = prevState.LastCollect.DiskC
		}
		if b.diskDTime.IsZero() {
			b.diskDTime = prevState.LastCollect.DiskD
		}
		if b.netCounterTime.IsZero() {
			b.netCounterTime = prevState.LastCollect.NetCounter
		}
	}

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
	b.mu.Lock()
	defer b.mu.Unlock()
	b.memTime = time
}

func (b *BasicStateSaver) SetDiskCMetric(time time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.diskCTime = time
}

func (b *BasicStateSaver) SetDiskDMetric(time time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.diskDTime = time
}

func (b *BasicStateSaver) SetNetCountersMetric(time time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.netCounterTime = time
}

func (b *BasicStateSaver) SetCpuMetric(time time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.cpuTime = time
}
