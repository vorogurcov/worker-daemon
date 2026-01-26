package job

import (
	"context"
	"fmt"
	"time"
)

type MonitoringJob struct {
	Name         string
	WorkTime     time.Duration
	WorkInterval time.Duration
	Callback     MonitoringCallback
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
				fmt.Printf("Finish monitoring job \"%v\"...\n", mj.Name)
				return nil
			}

			fmt.Printf("Stopping monitoring job \"%v\"...\n", mj.Name)
			return nil
		case t := <-ticker.C:
			if res, err := mj.Callback(mjCtx); err != nil {
				fmt.Printf("[MONITORING][%v][ERROR] %v\n", t.Format(time.RFC3339), err)
			} else {
				fmt.Printf("[MONITORING][%v] %s\n", t.Format(time.RFC3339), res.String())
			}
		}
	}
}
