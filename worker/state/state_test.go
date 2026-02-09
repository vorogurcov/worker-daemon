package state

import (
	"testing"
	"time"
)

func TestBasicStateSaver_Race(t *testing.T) {
	s := &BasicStateSaver{}

	go func() {
		for {
			s.SetCpuMetric(time.Now())
		}
	}()

	go func() {
		for {
			_ = s.GetShutdownState(true)
		}
	}()

	time.Sleep(100 * time.Millisecond)
}
