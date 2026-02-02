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

func (mj *MonitoringJob) Do(ctx context.Context) <-chan Result {
	resCh := make(chan Result)

	go func() {
		defer close(resCh)

		mjCtx, cancel := context.WithTimeout(ctx, mj.WorkTime)
		defer cancel()

		ticker := time.NewTicker(mj.WorkInterval)
		defer ticker.Stop()
		for {
			select {
			case <-mjCtx.Done():
				// WithCancelCause
				if ctx.Err() == context.DeadlineExceeded {
					resCh <- Result{
						Value: "Stopped by Worker (MaxWorkTime exceeded\n",
						Error: nil,
					}
					return
				}
				if mjCtx.Err() == context.DeadlineExceeded {
					resCh <- Result{
						Value: fmt.Sprintf("Finish monitoring job \"%v\"...\n", mj.Name),
						Error: nil,
					}
					return
				}
				resCh <- Result{
					Value: fmt.Sprintf("Stopping monitoring job \"%v\"...\n", mj.Name),
					Error: nil,
				}
				return
			case t := <-ticker.C:

				//ACTION
				if res, err := mj.Callback(mjCtx); err != nil {
					resCh <- Result{
						Value: fmt.Sprintf("[MONITORING][%v][ERROR] %v\n", t.Format(time.RFC3339), err),
						Error: nil,
					}
				} else {
					resCh <- Result{
						Value: fmt.Sprintf("[MONITORING][%v] %s\n", t.Format(time.RFC3339), res.String()),
						Error: nil,
					}
				}
			}
		}
	}()

	return resCh
}
