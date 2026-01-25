package job

import (
	"context"
	"fmt"
	"time"

	cpu2 "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type MonitoringJob struct {
	WorkTime     time.Duration
	WorkInterval time.Duration
	Callback     MonitoringStatisticsCallback
}

func (mj *MonitoringJob) Do(ctx context.Context) error {
	mjCtx, cancel := context.WithTimeout(ctx, mj.WorkTime)
	defer cancel()

	ticker := time.NewTicker(mj.WorkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-mjCtx.Done():
			// WithCancelCause
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("Stopped by Worker (MaxWorkTime exceeded)")
				return nil
			}

			if mjCtx.Err() == context.DeadlineExceeded {
				fmt.Println("Finish monitoring jobs...")
				return nil
			}

			fmt.Println("Stopping monitoring job...")
			return nil
		case t := <-ticker.C:
			//TODO: Как-то вынести каждый конкретный вызов в main, чтобы тут просто вызывать колбек
			// Мб типизировать с помощью interface MonitoringFunctionCallback
			v, _ := mem.VirtualMemory()
			cpu, _ := cpu2.Percent(0, false)
			d, _ := disk.Usage("C:")
			counters, _ := net.IOCounters(false)

			fmt.Printf("[MONITORING][%v] memory used: %.2f%%\n", t, v.UsedPercent)
			fmt.Printf("[MONITORING][%v] cpu used: %.2f%%\n", t, cpu[0])
			fmt.Printf("[MONITORING][%v] disk used: %.2f%%\n", t, d.UsedPercent)
			fmt.Printf("[MONITORING][%v] total bytes send: %v MiB\n", t,
				counters[0].BytesSent/1024/1024)

		}
	}
}
