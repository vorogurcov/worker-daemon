package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"main/job"
	"main/worker/state"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type BasicWorker struct {
	basicStateSaver *state.BasicStateSaver
	workerJobs      chan job.Job
	MaxWorkTime     time.Duration
	QueueSize       int
	wg              sync.WaitGroup
}

func NewWorker(stateSaver *state.BasicStateSaver, maxTime time.Duration, queueSize int) *BasicWorker {
	return &BasicWorker{
		workerJobs:      make(chan job.Job, queueSize),
		QueueSize:       queueSize,
		MaxWorkTime:     maxTime,
		basicStateSaver: stateSaver,
	}
}

func (bw *BasicWorker) ExecuteJobs(ctx context.Context) <-chan job.Result {
	resCh := make(chan job.Result, 10*bw.QueueSize)

	go func() {
		var workerCtx context.Context
		var cancel context.CancelFunc

		if bw.MaxWorkTime > 0 {
			workerCtx, cancel = context.WithTimeout(ctx, bw.MaxWorkTime)
		} else {
			workerCtx, cancel = context.WithCancel(ctx)
		}

		defer cancel()
		defer close(resCh)

		go func() {
			<-workerCtx.Done()
			if err := bw.Stop(); err != nil {
				resCh <- job.Result{
					Value: nil,
					Error: err,
				}
			}
			close(bw.workerJobs)
		}()

		for j := range bw.workerJobs {
			bw.wg.Add(1)
			go func(jobToRun job.Job) {
				defer bw.wg.Done()
				for res := range jobToRun.Do(workerCtx) {
					resCh <- res
				}
			}(j)
		}

		bw.wg.Wait()

	}()

	return resCh
}

func (bw *BasicWorker) AppendToJobs(ctx context.Context, job job.Job) {
	select {
	case <-ctx.Done():
		return
	default:
		bw.workerJobs <- job
		return
	}
}
func (bw *BasicWorker) Stop() error {
	shState := bw.basicStateSaver.GetShutdownState(true)

	const dir = "saves"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	lastPath := filepath.Join(dir, "last-save.json")

	if _, err := os.Stat(lastPath); err == nil {
		oldFile, err := os.Open(lastPath)
		if err == nil {
			var prev state.ShutdownState
			if err := json.NewDecoder(oldFile).Decode(&prev); err == nil {
				ts := prev.ShutdownTimeUTC.UTC().UnixMilli()
				historyName := filepath.Join(dir, fmt.Sprintf("save-%d.json", ts))
				oldFile.Close()
				_ = os.Rename(lastPath, historyName)
			} else {
				oldFile.Close()
			}
		}
	}

	tmpPath := lastPath + ".tmp"

	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	if err := enc.Encode(shState); err != nil {
		f.Close()
		return err
	}

	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, lastPath); err != nil {
		return err
	}

	return nil
}
