package job

import (
	"context"
	"fmt"
	"time"
)

type WaitingJob struct {
	WorkInterval time.Duration
	WorkTime     time.Duration
}

func (wj *WaitingJob) Do(ctx context.Context) <-chan Result {
	resCh := make(chan Result)

	go func() {
		defer close(resCh)
		// context PER CALL
		jobCtx, cancel := context.WithTimeout(ctx, wj.WorkTime)
		defer cancel()

		ticker := time.NewTicker(wj.WorkInterval)
		defer ticker.Stop()

		startTime := time.Now()

		for {
			select {
			case <-jobCtx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					resCh <- Result{
						Value: "Stopped by Worker (MaxWorkTime exceeded)\n",
						Error: nil,
					}
					return
				}
				if jobCtx.Err() == context.DeadlineExceeded {
					resCh <- Result{
						Value: "Fuck you! I'm done!\n",
						Error: nil,
					}
					return
				}

				resCh <- Result{
					Value: "Where do you go? Don't leave me!\n",
					Error: nil,
				}
				return
			case t := <-ticker.C:
				diff := t.Sub(startTime)

				resCh <- Result{
					Value: fmt.Sprintf("I'm still waiting for you... It's been %veconds!\n", diff),
					Error: nil,
				}
			}
		}
	}()

	return resCh
}
