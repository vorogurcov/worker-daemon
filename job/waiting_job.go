package job

import (
	"context"
	"fmt"
	"time"
)

type WaitingJob struct {
	WaitTime time.Duration
	WorkTime time.Duration
}

func (wj *WaitingJob) Do(ctx context.Context) error {
	// context PER CALL
	jobCtx, cancel := context.WithTimeout(ctx, wj.WorkTime)
	defer cancel()

	if v, ok := ctx.(Job); ok && v != nil {

	}

	ticker := time.NewTicker(wj.WaitTime)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-jobCtx.Done():

			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("Stopped by Worker (MaxWorkTime exceeded)")
				return nil
			}

			if jobCtx.Err() == context.DeadlineExceeded {
				fmt.Println("Fuck you! I'm done!")
				return nil
			}

			fmt.Println("Where do you go? Don't leave me!")
			return nil
		case t := <-ticker.C:
			diff := t.Sub(startTime)
			fmt.Printf("I'm still waiting for you... It's been %veconds!\n", diff)
		}
	}

}
