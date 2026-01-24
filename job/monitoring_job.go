package job

import (
	"context"
	"fmt"
	"time"
)

type MonitoringJob struct {
	WorkTime     time.Duration
	WorkInterval time.Duration
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
			//TODO: Do some monitoring act
			fmt.Printf("[MONITORING][%v] stat\n", t)
		}
	}
}
